package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
)

type ProviderRepository struct {
	Base Repository
}

func (w ProviderRepository) ListProviders(ctx context.Context, workspaceID int64) ([]sqlc.Provider, error) {
	return w.Base.DB.Queries.ListProvidersPostgresql(ctx, workspaceID)
}

func (w ProviderRepository) GetProvider(ctx context.Context, providerID int64) (sqlc.Provider, error) {
	return w.Base.DB.Queries.GetProviderPostgresql(ctx, providerID)
}

func (w ProviderRepository) UpdateSecurityPolicy(ctx context.Context, workspaceID, providerID int64, policy string) error {
	result, err := w.Base.DB.DB.ExecContext(ctx, `UPDATE providers SET security_policy = $3 WHERE id = $1 AND workspace_id = $2`, providerID, workspaceID, policy)
	if err != nil {
		return err
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (w ProviderRepository) UpdateProvider(ctx context.Context, providerID int64, name string, priority, weight int64) error {
	return w.Base.DB.Queries.UpdateProviderPostgresql(ctx, sqlc.UpdateProviderPostgresqlParams{
		ID:       providerID,
		Name:     name,
		Priority: priority,
		Weight:   weight,
	})
}

func (w ProviderRepository) GetProviderByName(ctx context.Context, name string) (sqlc.Provider, error) {
	return w.Base.DB.Queries.GetProviderByNamePostgresql(ctx, name)
}

func (w ProviderRepository) GetSmsConfig(ctx context.Context, providerID int64) (sqlc.SmsConfig, error) {
	return w.Base.DB.Queries.GetSmsConfigPostgresql(ctx, providerID)
}

func (w ProviderRepository) GetEmailConfig(ctx context.Context, providerID int64) (sqlc.EmailConfig, error) {
	return w.Base.DB.Queries.GetEmailConfigPostgresql(ctx, providerID)
}

func (w ProviderRepository) GetTelegramConfig(ctx context.Context, providerID int64) (sqlc.TelegramConfig, error) {
	return w.Base.DB.Queries.GetTelegramConfigPostgresql(ctx, providerID)
}

func (w ProviderRepository) GetSlackConfig(ctx context.Context, providerID int64) (sqlc.SlackConfig, error) {
	return w.Base.DB.Queries.GetSlackConfigPostgresql(ctx, providerID)
}

func (w ProviderRepository) GetWebhookConfig(ctx context.Context, providerID int64) (sqlc.WebhookConfig, error) {
	return w.Base.DB.Queries.GetWebhookConfigPostgresql(ctx, providerID)
}

func (w ProviderRepository) UpdateSmsConfig(ctx context.Context, providerID int64, url string, method int32, authTokenKey, authTokenValue *string, isAuthInHeader *bool, bodySchema *string, bodyEncoding, authType string, authUsername, accountSID, sender *string) error {
	return w.Base.DB.Queries.UpdateSmsConfigPostgresql(ctx, sqlc.UpdateSmsConfigPostgresqlParams{
		ProviderID:     providerID,
		Url:            url,
		Method:         int64(method),
		AuthTokenKey:   authTokenKey,
		Column5:        authTokenValue,
		IsAuthInHeader: isAuthInHeader,
		BodySchema:     bodySchema,
		BodyEncoding:   bodyEncoding,
		AuthType:       authType,
		AuthUsername:   authUsername,
		AccountSid:     accountSID,
		Sender:         sender,
	})
}

func (w ProviderRepository) UpdateEmailConfig(ctx context.Context, providerID int64, smtpHost string, smtpPort int64, username *string, password *string, fromEmail string) error {
	if username != nil && *username == "" {
		username = nil
	}
	if password != nil && *password == "" {
		password = nil
	}

	return w.Base.DB.Queries.UpdateEmailConfigPostgresql(ctx, sqlc.UpdateEmailConfigPostgresqlParams{
		ProviderID: providerID,
		SmtpHost:   smtpHost,
		SmtpPort:   smtpPort,
		Username:   username,
		Column4:    password,
		FromEmail:  fromEmail,
	})
}

func (w ProviderRepository) UpdateTelegramConfig(ctx context.Context, providerID int64, botToken string) error {
	return w.Base.DB.Queries.UpdateTelegramConfigPostgresql(ctx, sqlc.UpdateTelegramConfigPostgresqlParams{
		ProviderID: providerID,
		Column1:    botToken,
	})
}

func (w ProviderRepository) UpdateSlackConfig(ctx context.Context, providerID int64, webhookUrl string) error {
	return w.Base.DB.Queries.UpdateSlackConfigPostgresql(ctx, sqlc.UpdateSlackConfigPostgresqlParams{
		ProviderID: providerID,
		WebhookUrl: webhookUrl,
	})
}

func (w ProviderRepository) UpdateWebhookConfig(ctx context.Context, providerID int64, url string) error {
	return w.Base.DB.Queries.UpdateWebhookConfigPostgresql(ctx, sqlc.UpdateWebhookConfigPostgresqlParams{
		ProviderID: providerID,
		Url:        url,
	})
}

func (w ProviderRepository) CreateSmsConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, url string, method int32, authTokenKey, authTokenValue *string, isAuthInHeader *bool, bodySchema *string, bodyEncoding, authType string, authUsername, accountSID, sender *string) error {
	return w.Base.DB.Queries.CreateSmsConfigPostgresql(ctx, sqlc.CreateSmsConfigPostgresqlParams{
		WorkspaceID:    workspaceID,
		ProviderID:     providerID,
		OwnerID:        ownerID,
		Url:            url,
		Method:         int64(method),
		AuthTokenKey:   authTokenKey,
		AuthTokenValue: authTokenValue,
		IsAuthInHeader: isAuthInHeader,
		BodySchema:     bodySchema,
		BodyEncoding:   bodyEncoding,
		AuthType:       authType,
		AuthUsername:   authUsername,
		AccountSid:     accountSID,
		Sender:         sender,
	})
}

func (w ProviderRepository) CreateEmailConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, smtpHost string, smtpPort int64, fromEmail string, username *string, password *string) error {
	if username != nil && *username == "" {
		username = nil
	}
	if password != nil && *password == "" {
		password = nil
	}

	return w.Base.DB.Queries.CreateEmailConfigPostgresql(ctx, sqlc.CreateEmailConfigPostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		OwnerID:     ownerID,
		SmtpHost:    smtpHost,
		SmtpPort:    smtpPort,
		Username:    username,
		Password:    password,
		FromEmail:   fromEmail,
	})
}

func (w ProviderRepository) CreateTelegramConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, botToken string) error {
	// انتخاب دیتابیس و انجام عملیات
	return w.Base.DB.Queries.CreateTelegramConfigPostgresql(ctx, sqlc.CreateTelegramConfigPostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		OwnerID:     ownerID,
		BotToken:    botToken,
	})
}

func (w ProviderRepository) CreateSlackConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, webhookUrl string) error {
	// انتخاب دیتابیس و انجام عملیات
	return w.Base.DB.Queries.CreateSlackConfigPostgresql(ctx, sqlc.CreateSlackConfigPostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		OwnerID:     ownerID,
		WebhookUrl:  webhookUrl,
	})
}

func (w ProviderRepository) CreateWebhookConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, url string) error {
	return w.Base.DB.Queries.CreateWebhookConfigPostgresql(ctx, sqlc.CreateWebhookConfigPostgresqlParams{
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		OwnerID:     ownerID,
		Url:         url,
	})
}

func (w ProviderRepository) CreateProvider(ctx context.Context, workspaceID int64, ownerID int64, configType int32, name string, priority, weight int64) (sqlc.Provider, error) {
	return w.Base.DB.Queries.CreateProviderPostgresql(ctx, sqlc.CreateProviderPostgresqlParams{
		WorkspaceID: workspaceID,
		OwnerID:     ownerID,
		ConfigType:  int64(configType),
		Name:        name,
		Priority:    priority,
		Weight:      weight,
	})
}

func (w ProviderRepository) DeleteProvider(ctx context.Context, id int64) error {
	return w.Base.DB.Queries.DeleteProviderPostgresql(ctx, id)
}

func (w ProviderRepository) DeleteSmsConfig(ctx context.Context, providerID int64) error {
	return w.Base.DB.Queries.DeleteSmsConfigPostgresql(ctx, providerID)
}

func NewProviderRepository(base Repository) ProviderRepository {
	return ProviderRepository{Base: base}
}
