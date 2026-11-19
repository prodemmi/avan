package notifiers

import (
	"avan/internal/config"
	"avan/internal/db/sqlc"
	internalredis "avan/internal/redis"
	"avan/internal/services"
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (n Notifier) Start() {
	if n.Config.RedisConsumerCount < 1 {
		n.Logger.Fatal("failed to setup consumer: the consumer count cannot be 0")
		return
	}

	ctx := context.Background()
	client := n.MessageService.Repo.Base.Redis.Client()
	stream := internalredis.MessageQueueStream
	group := internalredis.MessageQueueGroup
	if err := client.XGroupCreateMkStream(ctx, stream, group, "0").Err(); err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		n.Logger.Fatal("failed to create Redis consumer group", zap.Error(err))
		return
	}

	consumer, err := os.Hostname()
	if err != nil {
		n.Logger.Fatal("failed to identify Redis consumer", zap.Error(err))
		return
	}

	claimCursor := "0-0"
	for {
		if err := n.MessageService.EnqueueDueScheduledMessages(ctx); err != nil {
			n.Logger.Error("failed to enqueue scheduled messages", zap.Error(err))
		}
		messages, next, err := client.XAutoClaim(ctx, &goredis.XAutoClaimArgs{
			Stream:   stream,
			Group:    group,
			Consumer: consumer,
			MinIdle:  time.Minute,
			Start:    claimCursor,
			Count:    10,
		}).Result()
		if err != nil {
			n.Logger.Fatal("failed to reclaim pending Redis messages", zap.Error(err))
			return
		}
		for _, message := range messages {
			n.consumeMessage(ctx, client, stream, group, message)
		}
		claimCursor = next
		if claimCursor != "0-0" {
			continue
		}

		streams, err := client.XReadGroup(ctx, &goredis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    time.Second,
		}).Result()
		if err == goredis.Nil {
			continue
		}
		if err != nil {
			n.Logger.Fatal("failed to read Redis queue", zap.Error(err))
			return
		}
		for _, result := range streams {
			for _, message := range result.Messages {
				n.consumeMessage(ctx, client, stream, group, message)
			}
		}
	}
}

func (n Notifier) consumeMessage(ctx context.Context, client *goredis.Client, stream, group string, queued goredis.XMessage) {
	payload, ok := queued.Values["payload"].(string)
	if !ok {
		n.Logger.Error("discarding malformed Redis queue message", zap.String("id", queued.ID))
		n.ackMessage(ctx, client, stream, group, queued.ID)
		return
	}

	var message sqlc.Message
	if err := json.Unmarshal([]byte(payload), &message); err != nil {
		n.Logger.Error("discarding invalid Redis queue message", zap.String("id", queued.ID), zap.Error(err))
		n.ackMessage(ctx, client, stream, group, queued.ID)
		return
	}

	provider, err := n.ProviderService.GetProvider(message.ProviderID)
	if err == nil {
		err = SendWithRetry(func() error { return n.SendMessage(provider, message) }, message.Retry)
	}
	if err != nil {
		n.Logger.Error("failed to send queued message", zap.Int64("message_id", message.ID), zap.Error(err))
		if markErr := n.MessageService.MarkMessageAsError(message.ID); markErr != nil {
			n.Logger.Error("failed to mark message as error", zap.Int64("message_id", message.ID), zap.Error(markErr))
			return
		}
		n.ackMessage(ctx, client, stream, group, queued.ID)
		return
	}
	if err := n.MessageService.MarkMessageAsSuccess(message.ID); err != nil {
		n.Logger.Error("failed to mark message as success", zap.Int64("message_id", message.ID), zap.Error(err))
		return
	}
	n.ackMessage(ctx, client, stream, group, queued.ID)
}

func (n Notifier) ackMessage(ctx context.Context, client *goredis.Client, stream, group, id string) {
	if err := client.XAck(ctx, stream, group, id).Err(); err != nil {
		n.Logger.Error("failed to acknowledge Redis queue message", zap.String("id", id), zap.Error(err))
		return
	}
	if err := client.XDel(ctx, stream, id).Err(); err != nil {
		n.Logger.Error("failed to remove acknowledged Redis queue message", zap.String("id", id), zap.Error(err))
	}
}

func NewNotifier(config config.Config, logger zap.Logger, providerService services.ProviderService, messageService services.MessageService) Notifier {
	return Notifier{Config: config, Logger: logger, ProviderService: providerService, MessageService: messageService}
}
