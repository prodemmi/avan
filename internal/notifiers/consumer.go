package notifiers

import (
	"avan/internal/config"
	"avan/internal/db/sqlc"
	"avan/internal/services"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Handler struct {
	Notifier Notifier
	Logger   zap.Logger
}

func (h Handler) Setup(session sarama.ConsumerGroupSession) error {

	return nil
}

func (h Handler) Cleanup(session sarama.ConsumerGroupSession) error {

	return nil
}

func (h Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message sqlc.Message

		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			h.Logger.Fatal("failed to unmarshal message value", zap.Error(err))
		}

		provider, err := h.Notifier.ProviderService.GetProvider(message.ProviderID)
		if err != nil {
			h.Logger.Error("failed to get provider info", zap.Error(err))
		}

		commitMessage := func(err error) {
			if err == nil {
				if err := h.Notifier.MessageService.MarkMessageAsSuccess(message.ID); err == nil {
					session.MarkMessage(msg, "")
				}
			} else {
				_ = h.Notifier.MessageService.MarkMessageAsError(message.ID)
			}
		}

		err = h.Notifier.SendMessage(provider, message)
		commitMessage(err)
	}
	return nil
}

func (n Notifier) Start() {
	consumersNumber := n.Config.KafkaConsumerCount

	if consumersNumber == 0 {
		n.Logger.Fatal("failed to setup consumer: the consumer count cannot be 0")
	}

	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.AutoCommit.Enable = false

	ctx := context.Background()
	consumer, err := sarama.NewConsumerGroup(n.Config.KafkaBrokers, "notifier-group-id", cfg)
	if err != nil {
		n.Logger.Fatal("failed to start consumer group", zap.Error(err))
	}

	topics := []string{n.Config.KafkaTopic}

	for range consumersNumber {
		err = consumer.Consume(ctx, topics, Handler{n, n.Logger})
		if err != nil {
			n.Logger.Fatal("failed to start consumer", zap.Error(err))
		}
	}
}

func NewNotifier(config config.Config, logger zap.Logger, providerService services.ProviderService, messageService services.MessageService) Notifier {
	return Notifier{Config: config, Logger: logger, ProviderService: providerService, MessageService: messageService}
}
