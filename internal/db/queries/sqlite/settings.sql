-- name: GetSettingsSqlite :many
SELECT *
from settings
WHERE workspace_id = ?
ORDER BY created_at DESC;

-- name: GetSettingSqlite :one
SELECT *
from settings
WHERE workspace_id = ? AND key = ?
LIMIT 1;

-- name: CreateSettingSqlite :exec
INSERT INTO settings (workspace_id, key, value) VALUES (?, ?, ?);

-- name: UpdateSettingSqlite :exec
UPDATE settings SET value = ? WHERE key = ? AND workspace_id = ?;