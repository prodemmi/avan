package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"time"
)

type MessageRepository struct {
	Base Repository
}

func (w MessageRepository) ListMessages(ctx context.Context, page int32, workspaceID int64, providerID *int64, templateID *int64, recipient *string, content *string, status *int16) ([]sqlc.Message, error) {
	return w.Base.DB.Queries.ListMessagesPostgresql(ctx, sqlc.ListMessagesPostgresqlParams{
		WorkspaceID: workspaceID,
		TemplateID:  templateID,
		ProviderID:  providerID,
		Recipient:   recipient,
		Content:     content,
		Status:      status,
	})
}

func (w MessageRepository) MessageCount(ctx context.Context, workspaceID int64, providerID *int64, templateID *int64, recipient *string, content *string, status *int16) (int64, error) {
	return w.Base.DB.Queries.MessageCountPostgresql(ctx, sqlc.MessageCountPostgresqlParams{
		WorkspaceID: workspaceID,
		TemplateID:  templateID,
		ProviderID:  providerID,
		Status:      status,
	})
}

func (w MessageRepository) CreateMessage(ctx context.Context, workspaceID, providerID int64, templateID *int64, recipient, content string, status int8, error *string, retry bool, scheduleAt, sentAt *time.Time, payload *string) (int64, error) {
	return w.Base.DB.Queries.CreateMessagePostgresql(ctx, sqlc.CreateMessagePostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		TemplateID:  templateID,
		Recipient:   recipient,
		Content:     content,
		Status:      int64(status),
		Error:       error,
		Retry:       retry,
		ScheduleAt:  scheduleAt,
		SentAt:      sentAt,
		Payload:     payload,
	})
}

func (w MessageRepository) UpdateMessageStatus(ctx context.Context, messageID int64, status int32) error {
	return w.Base.DB.Queries.UpdateMessageStatusPostgresql(ctx, sqlc.UpdateMessageStatusPostgresqlParams{
		ID:     messageID,
		Status: int64(status),
	})
}

func (w MessageRepository) UpdateMessageSentAt(ctx context.Context, messageID int64, sentAt *time.Time) error {
	return w.Base.DB.Queries.UpdateMessageSentAtPostgresql(ctx, sqlc.UpdateMessageSentAtPostgresqlParams{
		ID:     messageID,
		SentAt: sentAt,
	})
}

func (w MessageRepository) GetMessageSummary(ctx context.Context, workspaceID int64) (sqlc.MessageSummaryPostgresqlRow, error) {
	return w.Base.DB.Queries.MessageSummaryPostgresql(ctx, workspaceID)
}

func (w MessageRepository) GetProviderStats(ctx context.Context, workspaceID, providerID int64) (sqlc.GetProviderStatsPostgresqlRow, error) {
	return w.Base.DB.Queries.GetProviderStatsPostgresql(ctx, sqlc.GetProviderStatsPostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
	})
}

func (w MessageRepository) GetMessage(ctx context.Context, messageID int64) (sqlc.Message, error) {
	return w.Base.DB.Queries.GetMessagePostgresql(ctx, messageID)
}

func (w MessageRepository) GetTimeDistribution(ctx context.Context, workspaceID, providerID int64, startDate, endDate *time.Time) ([]sqlc.GetTimeDistributionPostgresqlRow, error) {
	return w.Base.DB.Queries.GetTimeDistributionPostgresql(ctx, sqlc.GetTimeDistributionPostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		StartDate:   startDate,
		EndDate:     endDate,
	})
}

func NewMessageRepository(base Repository) MessageRepository {
	return MessageRepository{Base: base}
}
