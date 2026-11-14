-- name: GetSettingsMysql :many
SELECT *
from settings
WHERE workspace_id = ?
ORDER BY created_at DESC;

-- name: GetSettingMysql :one
SELECT *
from settings
WHERE workspace_id = ? AND `key` = ?
LIMIT 1;

-- name: CreateSettingMysql :exec
INSERT INTO settings (workspace_id, `key`, value) VALUES (?, ?, ?);

-- name: UpdateSettingMysql :exec
UPDATE settings SET value = ? WHERE `key` = ? AND workspace_id = ?;