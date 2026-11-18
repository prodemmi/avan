-- name: GetApiKeysSqlite :many
SELECT *
from api_keys
WHERE workspace_id = ? AND owner_id = ?
ORDER BY created_at DESC;

-- name: CreateAPIKeySqlite :exec
INSERT INTO api_keys (workspace_id, owner_id, name, api_key, expire_at)
VALUES (?, ?, ?, ?, ?);

-- name: FindAPIKeySqlite :one
SELECT *
from api_keys
WHERE api_key = ?
LIMIT 1;

-- name: DeleteAPIKeySqlite :exec
DELETE FROM api_keys WHERE id = ?;
