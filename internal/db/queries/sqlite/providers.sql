-- name: ListProvidersSqlite :many
SELECT *
FROM providers
WHERE workspace_id = ?
  AND owner_id = ?
ORDER BY created_at DESC;

-- name: GetSmsConfigSqlite :one
SELECT *
FROM sms_configs
WHERE provider_id = ?;

-- name: GetEmailConfigSqlite :one
SELECT *
FROM email_configs
WHERE provider_id = ?;

-- name: GetTelegramConfigSqlite :one
SELECT *
FROM telegram_configs
WHERE provider_id = ?;

-- name: GetSlackConfigSqlite :one
SELECT *
FROM slack_configs
WHERE provider_id = ?;

-- name: GetWebhookConfigSqlite :one
SELECT *
FROM webhook_configs
WHERE provider_id = ?;

-- name: CreateSmsConfigSqlite :exec
INSERT INTO sms_configs (workspace_id, provider_id, owner_id, url, method, auth_token_key, auth_token_value,
                         is_auth_in_header)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateEmailConfigSqlite :exec
INSERT INTO email_configs (workspace_id, provider_id, owner_id, smtp_host, smtp_port, username, password, from_email)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateTelegramConfigSqlite :exec
INSERT INTO telegram_configs (workspace_id, provider_id, owner_id, bot_token)
VALUES (?, ?, ?, ?);

-- name: CreateSlackConfigSqlite :exec
INSERT INTO slack_configs (workspace_id, provider_id, owner_id, webhook_url)
VALUES (?, ?, ?, ?);

-- name: CreateWebhookConfigSqlite :exec
INSERT INTO webhook_configs (workspace_id, provider_id, owner_id, url)
VALUES (?, ?, ?, ?);

-- name: UpdateProviderSqlite :exec
UPDATE providers
SET name = ?
WHERE id = ?;

-- name: UpdateSmsConfigSqlite :exec
UPDATE sms_configs
SET url               = ?,
    method            = ?,
    auth_token_key    = ?,
    auth_token_value  = ?,
    is_auth_in_header = ?
WHERE provider_id = ?;

-- name: UpdateEmailConfigSqlite :exec
UPDATE email_configs
SET smtp_host  = ?,
    smtp_port  = ?,
    username   = ?,
    password   = ?,
    from_email = ?
WHERE provider_id = ?;

-- name: UpdateTelegramConfigSqlite :exec
UPDATE telegram_configs
SET bot_token = ?
WHERE provider_id = ?;

-- name: UpdateSlackConfigSqlite :exec
UPDATE slack_configs
SET webhook_url = ?
WHERE provider_id = ?;

-- name: UpdateWebhookConfigSqlite :exec
UPDATE webhook_configs
SET url = ?
WHERE provider_id = ?;

-- name: CreateProviderSqlite :one
INSERT INTO providers (workspace_id, owner_id, config_type, name)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetProviderSqlite :one
SELECT *
FROM providers
WHERE id = ?;

-- name: GetProviderByNameSqlite :one
SELECT *
FROM providers
WHERE name = ?;

-- name: DeleteProviderSqlite :exec
DELETE
FROM providers
WHERE id = ?;

-- name: DeleteSmsConfigSqlite :exec
DELETE
FROM sms_configs
WHERE provider_id = ?;
