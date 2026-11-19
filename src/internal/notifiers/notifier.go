package notifiers

import (
	"avan/internal/config"
	"avan/internal/db/sqlc"
	"avan/internal/services"
	"go.uber.org/zap"
	"time"
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

func SendWithRetry(send func() error, retry bool) error {
	attempts := 1
	if retry {
		attempts = 3
	}
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		if err = send(); err == nil {
			return nil
		}
		if attempt+1 < attempts {
			time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
		}
	}
	return err
}
