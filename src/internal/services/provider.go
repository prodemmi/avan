package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/redis"
	"avan/internal/repositories"
	"context"
	"time"
)

type ProviderService struct {
	Repo             repositories.ProviderRepository
	ProviderCache    redis.Cache[sqlc.Provider]
	SmsConfigCache   redis.Cache[sqlc.SmsConfig]
	EmailConfigCache redis.Cache[sqlc.EmailConfig]
	TelegramCache    redis.Cache[sqlc.TelegramConfig]
	SlackCache       redis.Cache[sqlc.SlackConfig]
	WebhookCache     redis.Cache[sqlc.WebhookConfig]
	WorkspaceService WorkspaceService
}

func (w ProviderService) GetProviders(workspaceID int64, ownerID int64) ([]sqlc.Provider, error) {
	ctx := context.Background()
	key := w.ProviderCache.GetKey("providers", workspaceID)
	return w.ProviderCache.GetOrSetArray(key, time.Hour*24, func() ([]sqlc.Provider, error) {
		return w.Repo.ListProviders(ctx, workspaceID)
	})
}

func (w ProviderService) GetProvider(providerID int64) (sqlc.Provider, error) {
	ctx := context.Background()
	key := w.ProviderCache.GetKey("provider", providerID)
	return w.ProviderCache.GetOrSet(key, time.Hour*24, func() (sqlc.Provider, error) {
		return w.Repo.GetProvider(ctx, providerID)
	})
}

func (w ProviderService) UpdateProvider(providerID int64, name string) error {
	ctx := context.Background()
	provider, err := w.GetProvider(providerID)
	if err != nil {
		return err
	}
	err = w.Repo.UpdateProvider(ctx, providerID, name)

	if err == nil {
		key := w.ProviderCache.GetStringKey("provider", providerID)
		w.ProviderCache.Delete(key...)
		w.invalidateProviders(provider.WorkspaceID)
	}

	return err
}

func (w ProviderService) GetProviderByName(name string) (sqlc.Provider, error) {
	ctx := context.Background()
	return w.Repo.GetProviderByName(ctx, name)
}

func (w ProviderService) GetSmsConfig(providerID int64) (sqlc.SmsConfig, error) {
	ctx := context.Background()
	key := w.SmsConfigCache.GetKey("sms-config", providerID)
	return w.SmsConfigCache.GetOrSet(key, time.Hour*24, func() (sqlc.SmsConfig, error) {
		return w.Repo.GetSmsConfig(ctx, providerID)
	})
}

func (w ProviderService) GetEmailConfig(providerID int64) (sqlc.EmailConfig, error) {
	ctx := context.Background()
	key := w.EmailConfigCache.GetKey("email-config", providerID)
	return w.EmailConfigCache.GetOrSet(key, time.Hour*24, func() (sqlc.EmailConfig, error) {
		return w.Repo.GetEmailConfig(ctx, providerID)
	})
}

func (w ProviderService) GetTelegramConfig(providerID int64) (sqlc.TelegramConfig, error) {
	ctx := context.Background()
	key := w.TelegramCache.GetKey("telegram-config", providerID)
	return w.TelegramCache.GetOrSet(key, time.Hour*24, func() (sqlc.TelegramConfig, error) {
		return w.Repo.GetTelegramConfig(ctx, providerID)
	})
}

func (w ProviderService) GetSlackConfig(providerID int64) (sqlc.SlackConfig, error) {
	ctx := context.Background()
	key := w.SlackCache.GetKey("slack-config", providerID)
	return w.SlackCache.GetOrSet(key, time.Hour*24, func() (sqlc.SlackConfig, error) {
		return w.Repo.GetSlackConfig(ctx, providerID)
	})
}

func (w ProviderService) GetWebhookConfig(providerID int64) (sqlc.WebhookConfig, error) {
	ctx := context.Background()
	key := w.WebhookCache.GetKey("webhook-config", providerID)
	return w.WebhookCache.GetOrSet(key, time.Hour*24, func() (sqlc.WebhookConfig, error) {
		return w.Repo.GetWebhookConfig(ctx, providerID)
	})
}

func (w ProviderService) UpdateSmsConfig(id int64, url string, method int32, authTokenKey, authTokenValue *string, isAuthInHeader *bool, bodySchema *string) error {
	ctx := context.Background()
	err := w.Repo.UpdateSmsConfig(ctx, id, url, method, authTokenKey, authTokenValue, isAuthInHeader, bodySchema)
	if err == nil {
		key := w.SmsConfigCache.GetStringKey("sms-config", id)
		w.SmsConfigCache.Delete(key...)
	}
	return err
}

func (w ProviderService) UpdateEmailConfig(id int64, smtpHost string, smtpPort int64, username *string, password *string, fromEmail string) error {
	ctx := context.Background()
	err := w.Repo.UpdateEmailConfig(ctx, id, smtpHost, smtpPort, username, password, fromEmail)
	if err == nil {
		key := w.EmailConfigCache.GetStringKey("email-config", id)
		w.EmailConfigCache.Delete(key...)
	}
	return err
}

func (w ProviderService) UpdateTelegramConfig(id int64, botToken string) error {
	ctx := context.Background()
	err := w.Repo.UpdateTelegramConfig(ctx, id, botToken)
	if err == nil {
		key := w.TelegramCache.GetStringKey("telegram-config", id)
		w.TelegramCache.Delete(key...)
	}
	return err
}

func (w ProviderService) UpdateSlackConfig(id int64, webhookUrl string) error {
	ctx := context.Background()
	err := w.Repo.UpdateSlackConfig(ctx, id, webhookUrl)
	if err == nil {
		key := w.SlackCache.GetStringKey("slack-config", id)
		w.SlackCache.Delete(key...)
	}
	return err
}

func (w ProviderService) UpdateWebhookConfig(id int64, url string) error {
	ctx := context.Background()
	err := w.Repo.UpdateWebhookConfig(ctx, id, url)
	if err == nil {
		key := w.WebhookCache.GetStringKey("webhook-config", id)
		w.WebhookCache.Delete(key...)
	}
	return err
}

func (w ProviderService) CreateSmsConfig(workspaceID int64, providerID int64, ownerID int64, url string, method int32, authTokenKey, authTokenValue *string, isAuthInHeader *bool, bodySchema *string) error {
	ctx := context.Background()
	return w.Repo.CreateSmsConfig(ctx, workspaceID, providerID, ownerID, url, method, authTokenKey, authTokenValue, isAuthInHeader, bodySchema)
}

func (w ProviderService) CreateEmailConfig(workspaceID int64, providerID int64, ownerID int64, smtpHost string, smtpPort int64, fromEmail string, username *string, password *string) error {
	ctx := context.Background()
	return w.Repo.CreateEmailConfig(ctx, workspaceID, providerID, ownerID, smtpHost, smtpPort, fromEmail, username, password)
}

func (w ProviderService) CreateTelegramConfig(workspaceID int64, providerID int64, ownerID int64, botToken string) error {
	ctx := context.Background()
	return w.Repo.CreateTelegramConfig(ctx, workspaceID, providerID, ownerID, botToken)
}

func (w ProviderService) CreateSlackConfig(workspaceID int64, providerID int64, ownerID int64, webhookUrl string) error {
	ctx := context.Background()
	return w.Repo.CreateSlackConfig(ctx, workspaceID, providerID, ownerID, webhookUrl)
}

func (w ProviderService) CreateWebhookConfig(workspaceID int64, providerID int64, ownerID int64, url string) error {
	ctx := context.Background()
	return w.Repo.CreateWebhookConfig(ctx, workspaceID, providerID, ownerID, url)
}

func (w ProviderService) CreateProvider(workspaceID int64, ownerID int64, configType int32, name string) (sqlc.Provider, error) {
	ctx := context.Background()
	provider, err := w.Repo.CreateProvider(ctx, workspaceID, ownerID, configType, name)
	if err == nil {
		w.invalidateProviders(workspaceID)
	}
	return provider, err
}

func (w ProviderService) DeleteProvider(id int64) error {
	ctx := context.Background()
	provider, err := w.GetProvider(id)
	if err != nil {
		return err
	}
	err = w.Repo.DeleteProvider(ctx, id)
	if err == nil {
		key := w.ProviderCache.GetStringKey("provider", id)
		w.ProviderCache.Delete(key...)
		w.invalidateProviders(provider.WorkspaceID)
	}
	return err
}

func (w ProviderService) invalidateProviders(workspaceID int64) {
	key := w.ProviderCache.GetStringKey("providers", workspaceID)
	_ = w.ProviderCache.Delete(key...)
}

func (w ProviderService) DeleteSmsConfig(id int64) error {
	ctx := context.Background()
	err := w.Repo.DeleteSmsConfig(ctx, id)
	if err == nil {
		key := w.SmsConfigCache.GetStringKey("sms-config", id)
		w.SmsConfigCache.Delete(key...)
	}
	return err
}

func (w ProviderService) CreateDefaultProviders() {
	for _, providerSettings := range w.Repo.Base.Config.Providers {
		var typeInt int32
		switch providerSettings.Type {
		case "sms":
			typeInt = 0
			break
		case "email":
			typeInt = 1
			break
		case "telegram":
			typeInt = 2
			break
		case "slack":
			typeInt = 3
			break
		case "webhook":
			typeInt = 4
			break
		}

		var methodInt int32
		switch providerSettings.Type {
		case "GET":
			methodInt = 0
			break
		case "get":
			methodInt = 0
			break
		case "POST":
			methodInt = 1
			break
		case "post":
			methodInt = 1
			break
		}

		if _, err := w.GetProviderByName(providerSettings.Name); err == nil {
			return
		}

		provider, _ := w.CreateProvider(1, 1, typeInt, providerSettings.Name)
		switch providerSettings.Type {
		case "sms":
			w.CreateSmsConfig(1, provider.ID, 1, providerSettings.Url, methodInt, &providerSettings.AuthTokenKey, &providerSettings.AuthTokenValue, &providerSettings.IsAuthInHeader, nil)
			break
		case "email":
			w.CreateEmailConfig(1, provider.ID, 1, providerSettings.SmtpHost, int64(providerSettings.SmtpPort), providerSettings.FromEmail, &providerSettings.Username, &providerSettings.Password)
			break
		case "telegram":
			w.CreateTelegramConfig(1, provider.ID, 1, providerSettings.BotToken)
			break
		case "slack":
			w.CreateSlackConfig(1, provider.ID, 1, providerSettings.WebhookUrl)
			break
		case "webhook":
			w.CreateWebhookConfig(1, provider.ID, 1, providerSettings.Url)
			break
		}
	}

}

func NewProviderService(repo repositories.ProviderRepository) ProviderService {
	providerCache := redis.NewCache[sqlc.Provider](repo.Base.Redis)
	smsCache := redis.NewCache[sqlc.SmsConfig](repo.Base.Redis)
	emailCache := redis.NewCache[sqlc.EmailConfig](repo.Base.Redis)
	telegramCache := redis.NewCache[sqlc.TelegramConfig](repo.Base.Redis)
	slackCache := redis.NewCache[sqlc.SlackConfig](repo.Base.Redis)
	webhookCache := redis.NewCache[sqlc.WebhookConfig](repo.Base.Redis)

	return ProviderService{
		Repo:             repo,
		ProviderCache:    providerCache,
		SmsConfigCache:   smsCache,
		EmailConfigCache: emailCache,
		TelegramCache:    telegramCache,
		SlackCache:       slackCache,
		WebhookCache:     webhookCache,
	}
}
