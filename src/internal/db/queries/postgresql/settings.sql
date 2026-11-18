-- name: GetSettingsPostgresql :many
SELECT *
from settings
WHERE workspace_id = $1
ORDER BY created_at DESC;

-- name: GetSettingPostgresql :one
SELECT *
from settings
WHERE workspace_id = $1 AND key = $2
LIMIT 1;

-- name: CreateSettingPostgresql :exec
INSERT INTO settings (workspace_id, key, value) VALUES ($1, $2, $3);

-- name: UpdateSettingPostgresql :exec
UPDATE settings SET value = $3 WHERE key = $2 AND workspace_id = $1;