package notifiers

import (
	"avan/internal/config"
	"avan/internal/db/sqlc"
	"avan/internal/services"
	"go.uber.org/zap"
)

type Notifier struct {
	Logger          zap.Logger
	Config          config.Config
	ProviderService services.ProviderService
	MessageService  services.MessageService
}

func (n Notifier) SendMessage(provider sqlc.Provider, message sqlc.Message) error {
	switch provider.ConfigType {
	case 0:
		return n.HandleSmsNotifier(provider, message)
	case 1:
		return n.HandleEmailNotifier(provider, message)
	case 2:
		return n.HandleTelegramNotifier(provider, message)
	case 3:
		return n.HandleSlackNotifier(provider, message)
	case 4:
		return n.HandleWebhookNotifier(provider, message)
	}
	return nil
}
