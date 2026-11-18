package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
	"time"
)

type MessageRepository struct {
	Base Repository
}

func (w MessageRepository) ListMessages(ctx context.Context, page int32, workspaceID int64, providerID *int64, templateID *int64, recipient *string, content *string, status *int16) ([]sqlc.Message, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		var nullTemplateID sql.NullInt64
		var nullProviderID sql.NullInt64
		var nullStatus sql.NullInt16

		if templateID != nil {
			nullTemplateID = sql.NullInt64{Int64: *templateID, Valid: true}
		} else {
			nullTemplateID = sql.NullInt64{Valid: false}
		}

		if providerID != nil {
			nullProviderID = sql.NullInt64{Int64: *providerID, Valid: true}
		} else {
			nullProviderID = sql.NullInt64{Valid: false}
		}

		if status != nil {
			nullStatus = sql.NullInt16{Int16: *status, Valid: true}
		} else {
			nullStatus = sql.NullInt16{Valid: false}
		}

		return w.Base.DB.Queries.ListMessagesMysql(ctx, sqlc.ListMessagesMysqlParams{
			WorkspaceID: workspaceID,
			TemplateID:  nullTemplateID,
			ProviderID:  nullProviderID,
			//Recipient:   recipient,
			//Content:     content,
			Status: nullStatus,
			Offset: page * 10,
		})
	case "postgres":
		return w.Base.DB.Queries.ListMessagesPostgresql(ctx, sqlc.ListMessagesPostgresqlParams{
			WorkspaceID: workspaceID,
			TemplateID:  templateID,
			ProviderID:  providerID,
			Recipient:   recipient,
			Content:     content,
			Status:      status,
		})
	case "sqlite":
		return w.Base.DB.Queries.ListMessagesSqlite(ctx, sqlc.ListMessagesSqliteParams{
			WorkspaceID: workspaceID,
			TemplateID:  templateID,
			ProviderID:  providerID,
			Recipient:   recipient,
			Content:     content,
			Status:      status,
		})
	}
	return nil, nil
}

func (w MessageRepository) MessageCount(ctx context.Context, workspaceID int64, providerID *int64, templateID *int64, recipient *string, content *string, status *int16) (int64, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		var nullTemplateID sql.NullInt64
		var nullProviderID sql.NullInt64
		var nullStatus sql.NullInt16

		if templateID != nil {
			nullTemplateID = sql.NullInt64{Int64: *templateID, Valid: true}
		} else {
			nullTemplateID = sql.NullInt64{Valid: false}
		}

		if providerID != nil {
			nullProviderID = sql.NullInt64{Int64: *providerID, Valid: true}
		} else {
			nullProviderID = sql.NullInt64{Valid: false}
		}

		if status != nil {
			nullStatus = sql.NullInt16{Int16: *status, Valid: true}
		} else {
			nullStatus = sql.NullInt16{Valid: false}
		}

		return w.Base.DB.Queries.MessageCountMysql(ctx, sqlc.MessageCountMysqlParams{
			WorkspaceID: workspaceID,
			TemplateID:  nullTemplateID,
			ProviderID:  nullProviderID,
			Status:      nullStatus,
		})
	case "postgres":
		return w.Base.DB.Queries.MessageCountPostgresql(ctx, sqlc.MessageCountPostgresqlParams{
			WorkspaceID: workspaceID,
			TemplateID:  templateID,
			ProviderID:  providerID,
			Status:      status,
		})
	case "sqlite":
		return w.Base.DB.Queries.MessageCountSqlite(ctx, sqlc.MessageCountSqliteParams{
			WorkspaceID: workspaceID,
			TemplateID:  templateID,
			ProviderID:  providerID,
			Status:      status,
		})
	}
	return 0, nil
}

func (w MessageRepository) CreateMessage(ctx context.Context, workspaceID, providerID int64, templateID *int64, recipient, content string, status int8, error *string, scheduleAt, sentAt *time.Time, payload *string) (int64, error) {
	var templateIDValue int64
	validTemplateID := templateID != nil
	if validTemplateID {
		templateIDValue = *templateID
	} else {
		templateIDValue = 0
	}

	var nullError sql.NullString
	if error != nil {
		nullError = sql.NullString{String: *error}
	} else {
		nullError = sql.NullString{Valid: false}
	}

	var nullPayload sql.NullString
	if payload != nil {
		nullPayload = sql.NullString{String: *payload}
	} else {
		nullPayload = sql.NullString{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateMessageMysql(ctx, sqlc.CreateMessageMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			TemplateID: sql.NullInt64{
				Int64: templateIDValue,
				Valid: validTemplateID,
			},
			Recipient: recipient,
			Content:   content,
			Status:    status,
			Error:     nullError,
			ScheduleAt: sql.NullTime{
				Time:  time.Now(),
				Valid: scheduleAt != nil,
			},
			SentAt: sql.NullTime{
				Time:  time.Now(),
				Valid: sentAt != nil,
			},
			Payload: nullPayload,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateMessagePostgresql(ctx, sqlc.CreateMessagePostgresqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			TemplateID: sql.NullInt64{
				Int64: templateIDValue,
				Valid: validTemplateID,
			},
			Recipient: recipient,
			Content:   content,
			Status:    int32(status),
			Error:     nullError,
			ScheduleAt: sql.NullTime{
				Time:  time.Now(),
				Valid: scheduleAt != nil,
			},
			SentAt: sql.NullTime{
				Time:  time.Now(),
				Valid: sentAt != nil,
			},
			Payload: nullPayload,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateMessageSqlite(ctx, sqlc.CreateMessageSqliteParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			TemplateID:  templateID,
			Recipient:   recipient,
			Content:     content,
			Status:      int64(status),
			Error:       error,
			ScheduleAt: sql.NullTime{
				Time:  time.Now(),
				Valid: scheduleAt != nil,
			},
			SentAt: sql.NullTime{
				Time:  time.Now(),
				Valid: sentAt != nil,
			},
			Payload: payload,
		})
	}
	return 0, nil
}

func (w MessageRepository) UpdateMessageStatus(ctx context.Context, messageID int64, status int32) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateMessageStatusMysql(ctx, sqlc.UpdateMessageStatusMysqlParams{
			Status: int8(status),
			ID:     messageID,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateMessageStatusPostgresql(ctx, sqlc.UpdateMessageStatusPostgresqlParams{
			ID:     int32(messageID),
			Status: status,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateMessageStatusSqlite(ctx, sqlc.UpdateMessageStatusSqliteParams{
			ID:     messageID,
			Status: int64(status),
		})
	}
	return nil
}

func (w MessageRepository) UpdateMessageSentAt(ctx context.Context, messageID int64, sentAt *time.Time) error {
	var nullSentAt sql.NullTime

	if sentAt != nil {
		nullSentAt = sql.NullTime{Time: *sentAt, Valid: true}
	} else {
		nullSentAt = sql.NullTime{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateMessageSentAtMysql(ctx, sqlc.UpdateMessageSentAtMysqlParams{
			ID:     messageID,
			SentAt: nullSentAt,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateMessageSentAtPostgresql(ctx, sqlc.UpdateMessageSentAtPostgresqlParams{
			ID:     int32(messageID),
			SentAt: nullSentAt,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateMessageSentAtSqlite(ctx, sqlc.UpdateMessageSentAtSqliteParams{
			ID:     messageID,
			SentAt: nullSentAt,
		})
	}
	return nil
}

func (w MessageRepository) GetMessageSummary(ctx context.Context, workspaceID int64) (sqlc.MessageSummaryMysqlRow, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.MessageSummaryMysql(ctx, workspaceID)
	case "postgres":
		return w.Base.DB.Queries.MessageSummaryMysql(ctx, workspaceID)
	case "sqlite":
		return w.Base.DB.Queries.MessageSummaryMysql(ctx, workspaceID)
	}
	return sqlc.MessageSummaryMysqlRow{}, nil
}

func (w MessageRepository) GetProviderStatsMysql(ctx context.Context, workspaceID, providerID int64) (sqlc.GetProviderStatsMysqlRow, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetProviderStatsMysql(ctx, sqlc.GetProviderStatsMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
		})
	case "postgres":
		return w.Base.DB.Queries.GetProviderStatsMysql(ctx, sqlc.GetProviderStatsMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
		})
	case "sqlite":
		return w.Base.DB.Queries.GetProviderStatsMysql(ctx, sqlc.GetProviderStatsMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
		})
	}
	return sqlc.GetProviderStatsMysqlRow{}, nil
}

func (w MessageRepository) GetMessage(ctx context.Context, messageID int64) (sqlc.Message, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetMessageMysql(ctx, messageID)
	case "postgres":
		return w.Base.DB.Queries.GetMessagePostgresql(ctx, int32(messageID))
	case "sqlite":
		return w.Base.DB.Queries.GetMessageSqlite(ctx, messageID)
	}
	return sqlc.Message{}, nil
}

func (w MessageRepository) GetTimeDistribution(ctx context.Context, workspaceID, providerID int64, startDate, endDate *time.Time) (sqlc.GetTimeDistributionMysqlRow, error) {
	var nullStartDate sql.NullTime
	var nullEndDate sql.NullTime

	if startDate != nil {
		nullStartDate = sql.NullTime{Time: *startDate, Valid: true}
	} else {
		nullStartDate = sql.NullTime{Valid: false}
	}

	if endDate != nil {
		nullEndDate = sql.NullTime{Time: *endDate, Valid: true}
	} else {
		nullEndDate = sql.NullTime{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetTimeDistributionMysql(ctx, sqlc.GetTimeDistributionMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			StartDate:   nullStartDate,
			EndDate:     nullEndDate,
		})
	case "postgres":
		return w.Base.DB.Queries.GetTimeDistributionMysql(ctx, sqlc.GetTimeDistributionMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			StartDate:   nullStartDate,
			EndDate:     nullEndDate,
		})
	case "sqlite":
		return w.Base.DB.Queries.GetTimeDistributionMysql(ctx, sqlc.GetTimeDistributionMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			StartDate:   nullStartDate,
			EndDate:     nullEndDate,
		})
	}
	return sqlc.GetTimeDistributionMysqlRow{}, nil
}

func NewMessageRepository(base Repository) MessageRepository {
	return MessageRepository{Base: base}
}
