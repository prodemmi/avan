package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestConcurrentIdempotencyClaimAndWorkspaceScope(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run the PostgreSQL idempotency integration test")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(20)
	ctx := context.Background()
	schema := "phase3_idempotency_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := db.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer db.ExecContext(ctx, "DROP SCHEMA "+schema+" CASCADE")
	if _, err := db.ExecContext(ctx, fmt.Sprintf(`CREATE TABLE %s.messages (
		id BIGSERIAL PRIMARY KEY,
		workspace_id BIGINT NOT NULL,
		provider_id BIGINT NOT NULL,
		template_id BIGINT,
		recipient VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		status INT NOT NULL,
		error TEXT,
		payload VARCHAR(255),
		schedule_at TIMESTAMP,
		sent_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		retry BOOLEAN NOT NULL DEFAULT FALSE,
		idempotency_key TEXT,
		request_hash TEXT
	)`, schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE UNIQUE INDEX ON %s.messages (workspace_id, idempotency_key)", schema)); err != nil {
		t.Fatal(err)
	}

	const concurrentRequests = 16
	connections := make([]*sql.Conn, concurrentRequests)
	for i := range connections {
		connections[i], err = db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := connections[i].ExecContext(ctx, "SET search_path TO "+schema); err != nil {
			t.Fatal(err)
		}
		defer connections[i].Close()
	}

	scheduleAt := time.Now().Add(time.Hour)
	params := sqlc.CreateIdempotentMessagePostgresqlParams{
		WorkspaceID:    11,
		ProviderID:     3,
		Recipient:      "+15555550100",
		Content:        "scheduled test",
		Status:         1,
		ScheduleAt:     &scheduleAt,
		IdempotencyKey: sql.NullString{String: "same-key", Valid: true},
		RequestHash:    sql.NullString{String: "same-request", Valid: true},
	}
	start := make(chan struct{})
	type result struct {
		message sqlc.Message
		created bool
		err     error
	}
	results := make(chan result, concurrentRequests)
	var group sync.WaitGroup
	for _, conn := range connections {
		group.Add(1)
		go func(conn *sql.Conn) {
			defer group.Done()
			<-start
			queries := sqlc.New(conn)
			message, err := queries.CreateIdempotentMessagePostgresql(ctx, params)
			if errors.Is(err, sql.ErrNoRows) {
				message, err = queries.GetMessageByIdempotencyKeyPostgresql(ctx, sqlc.GetMessageByIdempotencyKeyPostgresqlParams{
					WorkspaceID:    params.WorkspaceID,
					IdempotencyKey: params.IdempotencyKey,
				})
				results <- result{message: message, err: err}
				return
			}
			results <- result{message: message, created: err == nil, err: err}
		}(conn)
	}
	close(start)
	group.Wait()
	close(results)

	var winnerID int64
	createdCount := 0
	for result := range results {
		if result.err != nil {
			t.Fatal(result.err)
		}
		if result.created {
			createdCount++
		}
		if winnerID == 0 {
			winnerID = result.message.ID
		} else if result.message.ID != winnerID {
			t.Fatalf("duplicate claims returned IDs %d and %d", winnerID, result.message.ID)
		}
		if result.message.ScheduleAt == nil {
			t.Fatal("scheduled idempotency claim lost its schedule time")
		}
	}
	if createdCount != 1 {
		t.Fatalf("created %d rows for concurrent same-key requests, want exactly 1", createdCount)
	}

	otherWorkspace := params
	otherWorkspace.WorkspaceID = 12
	otherWorkspace.IdempotencyKey = sql.NullString{String: "same-key", Valid: true}
	other, err := sqlc.New(connections[0]).CreateIdempotentMessagePostgresql(ctx, otherWorkspace)
	if err != nil || other.ID == winnerID {
		t.Fatalf("same key must be independent across workspaces: message=%d err=%v", other.ID, err)
	}
	owned, err := sqlc.New(connections[0]).GetMessageStatusForWorkspacePostgresql(ctx, sqlc.GetMessageStatusForWorkspacePostgresqlParams{ID: winnerID, WorkspaceID: 11})
	if err != nil || owned.ID != winnerID {
		t.Fatalf("workspace inquiry should return the owned message: message=%d err=%v", owned.ID, err)
	}
	if _, err := sqlc.New(connections[0]).GetMessageStatusForWorkspacePostgresql(ctx, sqlc.GetMessageStatusForWorkspacePostgresqlParams{ID: winnerID, WorkspaceID: 12}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("cross-workspace inquiry error = %v, want sql.ErrNoRows", err)
	}
}
