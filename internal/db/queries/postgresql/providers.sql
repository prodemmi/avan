-- name: ListProvidersPostgresql :many
SELECT * FROM providers
WHERE workspace_id = $1 AND owner_id = $2
ORDER BY created_at DESC;

-- name: GetSmsConfigPostgresql :one
SELECT * FROM sms_configs
WHERE provider_id = $1;

-- name: GetEmailConfigPostgresql :one
SELECT * FROM email_configs
WHERE provider_id = $1;

-- name: GetTelegramConfigPostgresql :one
SELECT * FROM telegram_configs
WHERE provider_id = $1;

-- name: GetSlackConfigPostgresql :one
SELECT * FROM slack_configs
WHERE provider_id = $1;

-- name: GetWebhookConfigPostgresql :one
SELECT * FROM webhook_configs
WHERE provider_id = $1;

-- name: CreateSmsConfigPostgresql :exec
INSERT INTO sms_configs (workspace_id, provider_id, owner_id, url, method, auth_token_key, auth_token_value, is_auth_in_header)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateEmailConfigPostgresql :exec
INSERT INTO email_configs (workspace_id, provider_id, owner_id, smtp_host, smtp_port, username, password, from_email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateTelegramConfigPostgresql :exec
INSERT INTO telegram_configs (workspace_id, provider_id, owner_id, bot_token)
VALUES ($1, $2, $3, $4);

-- name: CreateSlackConfigPostgresql :exec
INSERT INTO slack_configs (workspace_id, provider_id, owner_id, webhook_url)
VALUES ($1, $2, $3, $4);

-- name: CreateWebhookConfigPostgresql :exec
INSERT INTO webhook_configs (workspace_id, provider_id, owner_id, url)
VALUES ($1, $2, $3, $4);

-- name: UpdateProviderPostgresql :exec
UPDATE providers
SET name = $2
WHERE id = $1;

-- name: UpdateSmsConfigPostgresql :exec
UPDATE sms_configs
SET url = $2,
    method = $3,
    auth_token_key = $4,
    auth_token_value = $5,
    is_auth_in_header  = $6
WHERE provider_id = $1;

-- name: UpdateEmailConfigPostgresql :exec
UPDATE email_configs
SET smtp_host = $1, smtp_port = $2, username = $3, password = $4, from_email = $5
WHERE provider_id = $6;

-- name: UpdateTelegramConfigPostgresql :exec
UPDATE telegram_configs
SET bot_token = $1
WHERE provider_id = $2;

-- name: UpdateSlackConfigPostgresql :exec
UPDATE slack_configs
SET webhook_url = $1
WHERE provider_id = $2;

-- name: UpdateWebhookConfigPostgresql :exec
UPDATE webhook_configs
SET url = $1
WHERE provider_id = $2;

-- name: CreateProviderPostgresql :one
INSERT INTO providers (workspace_id, owner_id, config_type, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProviderPostgresql :one
SELECT * FROM providers
WHERE id = $1;

-- name: GetProviderByNamePostgresql :one
SELECT * FROM providers
WHERE name = $1;

-- name: DeleteProviderPostgresql :exec
DELETE FROM providers
WHERE id = $1;

-- name: DeleteSmsConfigPostgresql :exec
DELETE FROM sms_configs
WHERE provider_id = $1;