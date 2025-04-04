-- name: ListProvidersMysql :many
SELECT *
FROM providers
WHERE workspace_id = ?
  AND owner_id = ?
ORDER BY created_at DESC;

-- name: GetSmsConfigMysql :one
SELECT *
FROM sms_configs
WHERE provider_id = ?;

-- name: GetEmailConfigMysql :one
SELECT *
FROM email_configs
WHERE provider_id = ?;

-- name: GetTelegramConfigMysql :one
SELECT *
FROM telegram_configs
WHERE provider_id = ?;

-- name: GetSlackConfigMysql :one
SELECT *
FROM slack_configs
WHERE provider_id = ?;

-- name: GetWebhookConfigMysql :one
SELECT *
FROM webhook_configs
WHERE provider_id = ?;

-- name: CreateSmsConfigMysql :exec
INSERT INTO sms_configs (workspace_id, provider_id, owner_id, url, method, auth_token_key, auth_token_value, is_auth_in_header)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateEmailConfigMysql :exec
INSERT INTO email_configs (workspace_id, provider_id, owner_id, smtp_host, smtp_port, username, password, from_email)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateTelegramConfigMysql :exec
INSERT INTO telegram_configs (workspace_id, provider_id, owner_id, bot_token)
VALUES (?, ?, ?, ?);

-- name: CreateSlackConfigMysql :exec
INSERT INTO slack_configs (workspace_id, provider_id, owner_id, webhook_url)
VALUES (?, ?, ?, ?);

-- name: CreateWebhookConfigMysql :exec
INSERT INTO webhook_configs (workspace_id, provider_id, owner_id, url)
VALUES (?, ?, ?, ?);

-- name: UpdateProviderMysql :exec
UPDATE providers
SET name = ?
WHERE id = ?;

-- name: UpdateSmsConfigMysql :exec
UPDATE sms_configs
SET url    = ?,
    method = ?,
    auth_token_key = ?,
    auth_token_value = ?,
    is_auth_in_header = ?
    WHERE provider_id = ?;

-- name: UpdateEmailConfigMysql :exec
UPDATE email_configs
SET smtp_host  = ?,
    smtp_port  = ?,
    username   = ?,
    password   = ?,
    from_email = ?
WHERE provider_id = ?;

-- name: UpdateTelegramConfigMysql :exec
UPDATE telegram_configs
SET bot_token = ?
WHERE provider_id = ?;

-- name: UpdateSlackConfigMysql :exec
UPDATE slack_configs
SET webhook_url = ?
WHERE provider_id = ?;

-- name: UpdateWebhookConfigMysql :exec
UPDATE webhook_configs
SET url = ?
WHERE provider_id = ?;

-- name: CreateProviderMysql :execlastid
INSERT INTO providers (workspace_id, owner_id, config_type, name)
VALUES (?, ?, ?, ?);

-- name: GetProviderMysql :one
SELECT *
FROM providers
WHERE id = ?;

-- name: GetProviderByNameMysql :one
SELECT *
FROM providers
WHERE name = ?;

-- name: DeleteProviderMysql :exec
DELETE
FROM providers
WHERE id = ?;

-- name: DeleteSmsConfigMysql :exec
DELETE
FROM sms_configs
WHERE provider_id = ?;
