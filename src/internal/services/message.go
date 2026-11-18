package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/dto"
	"avan/internal/redis"
	"avan/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type MessageService struct {
	Repo                         repositories.MessageRepository
	MessageCache                 redis.Cache[sqlc.Message]
	MessageSummaryCache          redis.Cache[sqlc.MessageSummaryPostgresqlRow]
	MessageProvidersStatsCache   redis.Cache[map[string]sqlc.GetProviderStatsPostgresqlRow]
	MessageTimeDistributionCache redis.Cache[map[string][]sqlc.GetTimeDistributionPostgresqlRow]
	MessageCountCache            redis.Cache[int64]
	WorkspaceService             WorkspaceService
	TemplateService              TemplateService
	ProviderService              ProviderService
}

func (w MessageService) GetMessages(page int32, workspaceID int64, providerID *int64, templateID *int64, recipient *string, content *string, status *int16) ([]sqlc.Message, error) {
	ctx := context.Background()
	key := w.MessageCache.GetKey("messages", workspaceID, providerID, templateID, status)
	return w.MessageCache.GetOrSetArray(key, time.Second*5, func() ([]sqlc.Message, error) {
		return w.Repo.ListMessages(ctx, page, workspaceID, providerID, templateID, recipient, content, status)
	})
}

func (w MessageService) GetMessage(messageID int64) (sqlc.Message, error) {
	ctx := context.Background()
	key := w.MessageCache.GetKey("message", messageID)
	return w.MessageCache.GetOrSet(key, time.Second*5, func() (sqlc.Message, error) {
		return w.Repo.GetMessage(ctx, messageID)
	})
}

func (w MessageService) MessageCount(workspaceID int64, providerID *int64, templateID *int64, recipient *string, content *string, status *int16) (int64, error) {
	ctx := context.Background()
	key := w.MessageCountCache.GetKey("messages", workspaceID, providerID, templateID, status)
	return w.MessageCountCache.GetOrSet(key, time.Second*5, func() (int64, error) {
		return w.Repo.MessageCount(ctx, workspaceID, providerID, templateID, recipient, content, status)
	})
}

func (w MessageService) GetMessageSummary(workspaceID int64) (sqlc.MessageSummaryPostgresqlRow, error) {
	ctx := context.Background()
	key := w.MessageSummaryCache.GetKey("messages", "summary", workspaceID)
	return w.MessageSummaryCache.GetOrSet(key, time.Second*5, func() (sqlc.MessageSummaryPostgresqlRow, error) {
		return w.Repo.GetMessageSummary(ctx, workspaceID)
	})
}

func (w MessageService) GetProvidersStats(workspaceID, ownerID int64) (map[string]sqlc.GetProviderStatsPostgresqlRow, error) {
	ctx := context.Background()
	key := w.MessageProvidersStatsCache.GetKey("messages", "provider-stats", workspaceID)
	return w.MessageProvidersStatsCache.GetOrSet(key, time.Second*5, func() (map[string]sqlc.GetProviderStatsPostgresqlRow, error) {
		providers, err := w.ProviderService.GetProviders(workspaceID, ownerID)
		if err != nil {
			return nil, err
		}
		data := make(map[string]sqlc.GetProviderStatsPostgresqlRow)

		for _, provider := range providers {
			result, err := w.Repo.GetProviderStats(ctx, workspaceID, provider.ID)
			if err != nil {
				return nil, err
			}
			data[provider.Name] = result
		}
		return data, nil
	})
}

func (w MessageService) GetTimeDistribution(workspaceID, ownerID int64, startDate, endDate *time.Time) (map[string][]sqlc.GetTimeDistributionPostgresqlRow, error) {
	ctx := context.Background()
	key := w.MessageTimeDistributionCache.GetKey("messages", "provider-distribution", workspaceID, startDate, endDate)
	return w.MessageTimeDistributionCache.GetOrSet(key, time.Second*5, func() (map[string][]sqlc.GetTimeDistributionPostgresqlRow, error) {
		providers, err := w.ProviderService.GetProviders(workspaceID, ownerID)
		if err != nil {
			return nil, err
		}
		data := make(map[string][]sqlc.GetTimeDistributionPostgresqlRow)

		for _, provider := range providers {
			result, err := w.Repo.GetTimeDistribution(ctx, workspaceID, provider.ID, startDate, endDate)
			if err != nil {
				return nil, err
			}
			data[provider.Name] = result
		}
		return data, nil
	})
}

func (w MessageService) CreateMessage(workspaceID, providerID int64, templateID *int64, recipient, content string, status int8, error *string, scheduleAt, sentAt *time.Time, payload *string) (int64, error) {
	ctx := context.Background()
	key := w.MessageCache.GetKey("messages", workspaceID, providerID, templateID, status)
	err := w.MessageCache.Delete(w.MessageCache.GetStringKey(key)...)
	if err != nil {
		return 0, err
	}
	return w.Repo.CreateMessage(ctx, workspaceID, providerID, templateID, recipient, content, status, error, scheduleAt, sentAt, payload)
}

func (w MessageService) ProduceMessage(workspaceID, providerID int64, templateID *int64, recipient, content string, scheduleAt *time.Time, payload *string) error {
	var status int8
	status = 1

	sentAt := time.Now()
	id, err := w.CreateMessage(workspaceID, providerID, templateID, recipient, content, status, nil, scheduleAt, &sentAt, payload)
	if err != nil {
		return err
	}

	now := time.Now()
	message := sqlc.Message{
		ID:          id,
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		TemplateID:  templateID,
		Recipient:   recipient,
		Content:     content,
		Status:      int64(status),
		ScheduleAt:  scheduleAt,
		Payload:     payload,
		SentAt:      &sentAt,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	return w.enqueueMessage(message)
}

func (w MessageService) ProduceMessageDirectly(providerID, messageID int64) error {
	message, err := w.GetMessage(messageID)
	if err != nil {
		return err
	}

	return w.enqueueMessage(message)
}

func (w MessageService) enqueueMessage(message sqlc.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return w.Repo.Base.Redis.Client().XAdd(context.Background(), &goredis.XAddArgs{
		Stream: redis.MessageQueueStream,
		Values: map[string]interface{}{"payload": string(payload)},
	}).Err()
}

type TemplateData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (w MessageService) BindMessageData(templateID int64, data *[]dto.TemplateData) (string, error) {
	template, err := w.TemplateService.GetTemplate(templateID)
	if err != nil {
		return "", err
	}
	content := template.Content

	if data == nil {
		return content, nil
	}

	for _, d := range *data {
		if strings.HasPrefix(d.Key, "$") {
			content = strings.ReplaceAll(content, d.Key, d.Value)
		} else {
			content = strings.ReplaceAll(content, fmt.Sprintf("$%s", d.Key), d.Value)
		}
	}
	return content, nil
}

func (w MessageService) UpdateMessageSentAt(messageID int64, sentAt *time.Time) error {
	ctx := context.Background()
	key := w.MessageCache.GetKey("message", messageID)
	err := w.MessageCache.Delete(w.MessageCache.GetStringKey(key)...)
	if err != nil {
		return err
	}
	return w.Repo.UpdateMessageSentAt(ctx, messageID, sentAt)
}

func (w MessageService) MarkMessageAsSuccess(messageID int64) error {
	ctx := context.Background()
	key := w.MessageCache.GetKey("message", messageID)
	err := w.MessageCache.Delete(w.MessageCache.GetStringKey(key)...)
	if err != nil {
		return err
	}
	err = w.Repo.UpdateMessageStatus(ctx, messageID, 0)
	if err != nil {
		return err
	}
	now := time.Now()
	return w.UpdateMessageSentAt(messageID, &now)
}

func (w MessageService) MarkMessageAsPending(messageID int64) error {
	ctx := context.Background()
	key := w.MessageCache.GetKey("message", messageID)
	err := w.MessageCache.Delete(w.MessageCache.GetStringKey(key)...)
	if err != nil {
		return err
	}
	return w.Repo.UpdateMessageStatus(ctx, messageID, 1)
}

func (w MessageService) MarkMessageAsError(messageID int64) error {
	ctx := context.Background()
	key := w.MessageCache.GetKey("message", messageID)
	err := w.MessageCache.Delete(w.MessageCache.GetStringKey(key)...)
	if err != nil {
		return err
	}
	return w.Repo.UpdateMessageStatus(ctx, messageID, 2)
}

func NewMessageService(repo repositories.MessageRepository, templateService TemplateService, providerService ProviderService) MessageService {
	messageCache := redis.NewCache[sqlc.Message](repo.Base.Redis)
	messageCountCache := redis.NewCache[int64](repo.Base.Redis)
	messageSummaryCache := redis.NewCache[sqlc.MessageSummaryPostgresqlRow](repo.Base.Redis)
	messageProvidersStatsCache := redis.NewCache[map[string]sqlc.GetProviderStatsPostgresqlRow](repo.Base.Redis)
	messageTimeDistributionCache := redis.NewCache[map[string][]sqlc.GetTimeDistributionPostgresqlRow](repo.Base.Redis)

	return MessageService{
		Repo:                         repo,
		MessageCache:                 messageCache,
		MessageCountCache:            messageCountCache,
		MessageSummaryCache:          messageSummaryCache,
		TemplateService:              templateService,
		ProviderService:              providerService,
		MessageProvidersStatsCache:   messageProvidersStatsCache,
		MessageTimeDistributionCache: messageTimeDistributionCache,
	}
}
