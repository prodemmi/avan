package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
	"fmt"
)

type ProviderRepository struct {
	Base Repository
}

func (w ProviderRepository) ListProviders(ctx context.Context, workspaceID int64, ownerID int64) ([]sqlc.Provider, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.ListProvidersMysql(ctx, sqlc.ListProvidersMysqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	case "postgres":
		return w.Base.DB.Queries.ListProvidersPostgresql(ctx, sqlc.ListProvidersPostgresqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	case "sqlite":
		return w.Base.DB.Queries.ListProvidersSqlite(ctx, sqlc.ListProvidersSqliteParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	}
	return nil, nil
}

func (w ProviderRepository) GetProvider(ctx context.Context, providerID int64) (sqlc.Provider, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetProviderMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.GetProviderPostgresql(ctx, int32(providerID))
	case "sqlite":
		return w.Base.DB.Queries.GetProviderSqlite(ctx, providerID)
	}
	return sqlc.Provider{}, nil
}

func (w ProviderRepository) UpdateProvider(ctx context.Context, providerID int64, name string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateProviderMysql(ctx, sqlc.UpdateProviderMysqlParams{
			ID:   providerID,
			Name: name,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateProviderPostgresql(ctx, sqlc.UpdateProviderPostgresqlParams{
			ID:   int32(providerID),
			Name: name,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateProviderSqlite(ctx, sqlc.UpdateProviderSqliteParams{
			ID:   providerID,
			Name: name,
		})
	}
	return nil
}

func (w ProviderRepository) GetProviderByName(ctx context.Context, name string) (sqlc.Provider, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetProviderByNameMysql(ctx, name)
	case "postgres":
		return w.Base.DB.Queries.GetProviderByNamePostgresql(ctx, name)
	case "sqlite":
		return w.Base.DB.Queries.GetProviderByNameSqlite(ctx, name)
	}
	return sqlc.Provider{}, nil
}

func (w ProviderRepository) GetSmsConfig(ctx context.Context, providerID int64) (sqlc.SmsConfig, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetSmsConfigMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.GetSmsConfigPostgresql(ctx, providerID)
	case "sqlite":
		return w.Base.DB.Queries.GetSmsConfigSqlite(ctx, providerID)
	}
	return sqlc.SmsConfig{}, nil
}

func (w ProviderRepository) GetEmailConfig(ctx context.Context, providerID int64) (sqlc.EmailConfig, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetEmailConfigMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.GetEmailConfigPostgresql(ctx, providerID)
	case "sqlite":
		return w.Base.DB.Queries.GetEmailConfigSqlite(ctx, providerID)
	}
	return sqlc.EmailConfig{}, nil
}

func (w ProviderRepository) GetTelegramConfig(ctx context.Context, providerID int64) (sqlc.TelegramConfig, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetTelegramConfigMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.GetTelegramConfigPostgresql(ctx, providerID)
	case "sqlite":
		return w.Base.DB.Queries.GetTelegramConfigSqlite(ctx, providerID)
	}
	return sqlc.TelegramConfig{}, nil
}

func (w ProviderRepository) GetSlackConfig(ctx context.Context, providerID int64) (sqlc.SlackConfig, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetSlackConfigMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.GetSlackConfigPostgresql(ctx, providerID)
	case "sqlite":
		return w.Base.DB.Queries.GetSlackConfigSqlite(ctx, providerID)
	}
	return sqlc.SlackConfig{}, nil
}

func (w ProviderRepository) GetWebhookConfig(ctx context.Context, providerID int64) (sqlc.WebhookConfig, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetWebhookConfigMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.GetWebhookConfigPostgresql(ctx, providerID)
	case "sqlite":
		return w.Base.DB.Queries.GetWebhookConfigSqlite(ctx, providerID)
	}
	return sqlc.WebhookConfig{}, nil
}

func (w ProviderRepository) UpdateSmsConfig(ctx context.Context, providerID int64, url string, method int32, authTokenKey, authTokenValue *string, isAuthInHeader *bool) error {
	var nullAuthTokenKey sql.NullString
	var nullAuthTokenValue sql.NullString
	var nullIsAuthInHeader sql.NullBool

	if authTokenKey != nil {
		nullAuthTokenKey = sql.NullString{String: *authTokenKey, Valid: true}
	} else {
		nullAuthTokenKey = sql.NullString{Valid: false}
	}
	if authTokenValue != nil {
		nullAuthTokenValue = sql.NullString{String: *authTokenValue, Valid: true}
	} else {
		nullAuthTokenValue = sql.NullString{Valid: false}
	}
	if isAuthInHeader != nil {
		nullIsAuthInHeader = sql.NullBool{Bool: *isAuthInHeader, Valid: true}
	} else {
		nullIsAuthInHeader = sql.NullBool{Valid: false}
	}

	fmt.Println("nullIsAuthInHeader", nullIsAuthInHeader)

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateSmsConfigMysql(ctx, sqlc.UpdateSmsConfigMysqlParams{
			Url:            url,
			Method:         method,
			ProviderID:     providerID,
			AuthTokenKey:   nullAuthTokenKey,
			AuthTokenValue: nullAuthTokenValue,
			IsAuthInHeader: nullIsAuthInHeader,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateSmsConfigPostgresql(ctx, sqlc.UpdateSmsConfigPostgresqlParams{
			ProviderID:     providerID,
			Url:            url,
			Method:         method,
			AuthTokenKey:   nullAuthTokenKey,
			AuthTokenValue: nullAuthTokenValue,
			IsAuthInHeader: nullIsAuthInHeader,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateSmsConfigSqlite(ctx, sqlc.UpdateSmsConfigSqliteParams{
			ProviderID:     providerID,
			Url:            url,
			Method:         int64(method),
			AuthTokenKey:   authTokenKey,
			AuthTokenValue: authTokenValue,
			IsAuthInHeader: isAuthInHeader,
		})
	}
	return nil
}

func (w ProviderRepository) UpdateEmailConfig(ctx context.Context, providerID int64, smtpHost string, smtpPort int64, username *string, password *string, fromEmail string) error {
	// بررسی درست بودن username و password
	var nullUsername sql.NullString
	var nullPassword sql.NullString

	if username != nil && *username != "" {
		nullUsername = sql.NullString{String: *username, Valid: true}
	} else {
		nullUsername = sql.NullString{Valid: false}
	}

	if password != nil && *password != "" {
		nullPassword = sql.NullString{String: *password, Valid: true}
	} else {
		nullPassword = sql.NullString{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateEmailConfigMysql(ctx, sqlc.UpdateEmailConfigMysqlParams{
			ProviderID: providerID,
			SmtpHost:   smtpHost,
			SmtpPort:   smtpPort,
			Username:   nullUsername,
			Password:   nullPassword,
			FromEmail:  fromEmail,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateEmailConfigPostgresql(ctx, sqlc.UpdateEmailConfigPostgresqlParams{
			ProviderID: providerID,
			SmtpHost:   smtpHost,
			SmtpPort:   smtpPort,
			Username:   nullUsername,
			Password:   nullPassword,
			FromEmail:  fromEmail,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateEmailConfigSqlite(ctx, sqlc.UpdateEmailConfigSqliteParams{
			ProviderID: providerID,
			SmtpHost:   smtpHost,
			SmtpPort:   smtpPort,
			Username:   username,
			Password:   password,
			FromEmail:  fromEmail,
		})
	}
	return nil
}

func (w ProviderRepository) UpdateTelegramConfig(ctx context.Context, providerID int64, botToken string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateTelegramConfigMysql(ctx, sqlc.UpdateTelegramConfigMysqlParams{
			ProviderID: providerID,
			BotToken:   botToken,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateTelegramConfigPostgresql(ctx, sqlc.UpdateTelegramConfigPostgresqlParams{
			ProviderID: providerID,
			BotToken:   botToken,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateTelegramConfigSqlite(ctx, sqlc.UpdateTelegramConfigSqliteParams{
			ProviderID: providerID,
			BotToken:   botToken,
		})
	}
	return nil
}

func (w ProviderRepository) UpdateSlackConfig(ctx context.Context, providerID int64, webhookUrl string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateSlackConfigMysql(ctx, sqlc.UpdateSlackConfigMysqlParams{
			ProviderID: providerID,
			WebhookUrl: webhookUrl,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateSlackConfigPostgresql(ctx, sqlc.UpdateSlackConfigPostgresqlParams{
			ProviderID: providerID,
			WebhookUrl: webhookUrl,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateSlackConfigSqlite(ctx, sqlc.UpdateSlackConfigSqliteParams{
			ProviderID: providerID,
			WebhookUrl: webhookUrl,
		})
	}
	return nil
}

func (w ProviderRepository) UpdateWebhookConfig(ctx context.Context, providerID int64, url string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateWebhookConfigMysql(ctx, sqlc.UpdateWebhookConfigMysqlParams{
			ProviderID: providerID,
			Url:        url,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateWebhookConfigPostgresql(ctx, sqlc.UpdateWebhookConfigPostgresqlParams{
			ProviderID: providerID,
			Url:        url,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateWebhookConfigSqlite(ctx, sqlc.UpdateWebhookConfigSqliteParams{
			ProviderID: providerID,
			Url:        url,
		})
	}
	return nil
}

func (w ProviderRepository) CreateSmsConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, url string, method int32, authTokenKey, authTokenValue *string, isAuthInHeader *bool) error {
	var nullAuthTokenKey sql.NullString
	var nullAuthTokenValue sql.NullString
	var nullIsAuthInHeader sql.NullBool

	if authTokenKey != nil {
		nullAuthTokenKey = sql.NullString{String: *authTokenKey}
	} else {
		nullAuthTokenKey = sql.NullString{Valid: false}
	}
	if authTokenValue != nil {
		nullAuthTokenValue = sql.NullString{String: *authTokenValue}
	} else {
		nullAuthTokenValue = sql.NullString{Valid: false}
	}
	if isAuthInHeader != nil {
		nullIsAuthInHeader = sql.NullBool{Bool: *isAuthInHeader}
	} else {
		nullIsAuthInHeader = sql.NullBool{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateSmsConfigMysql(ctx, sqlc.CreateSmsConfigMysqlParams{
			WorkspaceID:    workspaceID,
			ProviderID:     providerID,
			OwnerID:        ownerID,
			Url:            url,
			Method:         method,
			AuthTokenKey:   nullAuthTokenKey,
			AuthTokenValue: nullAuthTokenValue,
			IsAuthInHeader: nullIsAuthInHeader,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateSmsConfigPostgresql(ctx, sqlc.CreateSmsConfigPostgresqlParams{
			WorkspaceID:    workspaceID,
			ProviderID:     providerID,
			OwnerID:        ownerID,
			Url:            url,
			Method:         method,
			AuthTokenKey:   nullAuthTokenKey,
			AuthTokenValue: nullAuthTokenValue,
			IsAuthInHeader: nullIsAuthInHeader,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateSmsConfigSqlite(ctx, sqlc.CreateSmsConfigSqliteParams{
			WorkspaceID:    workspaceID,
			ProviderID:     providerID,
			OwnerID:        ownerID,
			Url:            url,
			Method:         int64(method),
			AuthTokenKey:   authTokenKey,
			AuthTokenValue: authTokenValue,
			IsAuthInHeader: isAuthInHeader,
		})
	}
	return nil
}

func (w ProviderRepository) CreateEmailConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, smtpHost string, smtpPort int64, fromEmail string, username *string, password *string) error {
	// بررسی درست بودن username و password
	validUserName := username != nil && *username != ""
	validUserPassword := password != nil && *password != ""

	// انتخاب دیتابیس و انجام عملیات
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateEmailConfigMysql(ctx, sqlc.CreateEmailConfigMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			SmtpHost:    smtpHost,
			SmtpPort:    smtpPort,
			Username: sql.NullString{
				String: *username, // استفاده از مقدار اشاره‌گر
				Valid:  validUserName,
			},
			Password: sql.NullString{
				String: *password, // استفاده از مقدار اشاره‌گر
				Valid:  validUserPassword,
			},
			FromEmail: fromEmail,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateEmailConfigPostgresql(ctx, sqlc.CreateEmailConfigPostgresqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			SmtpHost:    smtpHost, // استفاده از مقدار اشاره‌گر
			SmtpPort:    smtpPort,
			Username: sql.NullString{
				String: *username, // استفاده از مقدار اشاره‌گر
				Valid:  validUserName,
			},
			Password: sql.NullString{
				String: *password, // استفاده از مقدار اشاره‌گر
				Valid:  validUserPassword,
			},
			FromEmail: fromEmail,
		})
	case "sqlite":
		// استفاده از اشاره‌گر برای ارسال username و password
		return w.Base.DB.Queries.CreateEmailConfigSqlite(ctx, sqlc.CreateEmailConfigSqliteParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			SmtpHost:    smtpHost, // استفاده از مقدار اشاره‌گر
			SmtpPort:    smtpPort,
			Username:    username, // ارسال مستقیم اشاره‌گر
			Password:    password, // ارسال مستقیم اشاره‌گر
			FromEmail:   fromEmail,
		})
	}
	return nil
}

func (w ProviderRepository) CreateTelegramConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, botToken string) error {
	// انتخاب دیتابیس و انجام عملیات
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateTelegramConfigMysql(ctx, sqlc.CreateTelegramConfigMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			BotToken:    botToken,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateTelegramConfigPostgresql(ctx, sqlc.CreateTelegramConfigPostgresqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			BotToken:    botToken,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateTelegramConfigSqlite(ctx, sqlc.CreateTelegramConfigSqliteParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			BotToken:    botToken,
		})
	}
	return nil
}

func (w ProviderRepository) CreateSlackConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, webhookUrl string) error {
	// انتخاب دیتابیس و انجام عملیات
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateSlackConfigMysql(ctx, sqlc.CreateSlackConfigMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			WebhookUrl:  webhookUrl,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateSlackConfigPostgresql(ctx, sqlc.CreateSlackConfigPostgresqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			WebhookUrl:  webhookUrl,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateSlackConfigSqlite(ctx, sqlc.CreateSlackConfigSqliteParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			WebhookUrl:  webhookUrl,
		})
	}
	return nil
}

func (w ProviderRepository) CreateWebhookConfig(ctx context.Context, workspaceID int64, providerID int64, ownerID int64, url string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateWebhookConfigMysql(ctx, sqlc.CreateWebhookConfigMysqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			Url:         url,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateWebhookConfigPostgresql(ctx, sqlc.CreateWebhookConfigPostgresqlParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			Url:         url,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateWebhookConfigSqlite(ctx, sqlc.CreateWebhookConfigSqliteParams{
			WorkspaceID: workspaceID,
			ProviderID:  providerID,
			OwnerID:     ownerID,
			Url:         url,
		})
	}
	return nil
}

func (w ProviderRepository) CreateProvider(ctx context.Context, workspaceID int64, ownerID int64, configType int32, name string) (sqlc.Provider, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		id, err := w.Base.DB.Queries.CreateProviderMysql(ctx, sqlc.CreateProviderMysqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			ConfigType:  configType,
			Name:        name,
		})
		if err != nil {
			return sqlc.Provider{}, err
		}
		provider, err := w.Base.DB.Queries.GetProviderMysql(ctx, id)
		return provider, nil
	case "postgres":
		return w.Base.DB.Queries.CreateProviderPostgresql(ctx, sqlc.CreateProviderPostgresqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			ConfigType:  configType,
			Name:        name,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateProviderSqlite(ctx, sqlc.CreateProviderSqliteParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			ConfigType:  int64(configType),
			Name:        name,
		})
	}
	return sqlc.Provider{}, nil
}

func (w ProviderRepository) DeleteProvider(ctx context.Context, id int64) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.DeleteProviderMysql(ctx, id)
	case "postgres":
		return w.Base.DB.Queries.DeleteProviderPostgresql(ctx, int32(id))
	case "sqlite":
		return w.Base.DB.Queries.DeleteProviderSqlite(ctx, id)
	}
	return nil
}

func (w ProviderRepository) DeleteSmsConfig(ctx context.Context, providerID int64) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.DeleteSmsConfigMysql(ctx, providerID)
	case "postgres":
		return w.Base.DB.Queries.DeleteProviderPostgresql(ctx, int32(providerID))
	case "sqlite":
		return w.Base.DB.Queries.DeleteProviderSqlite(ctx, providerID)
	}
	return nil
}

func NewProviderRepository(base Repository) ProviderRepository {
	return ProviderRepository{Base: base}
}
