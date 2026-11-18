-- name: GetApiKeysPostgresql :many
SELECT *
from api_keys
WHERE workspace_id = $1
ORDER BY created_at DESC;

-- name: CreateAPIKeyPostgresql :exec
INSERT INTO api_keys (workspace_id, owner_id, name, api_key, expire_at)
VALUES ($1, $2, $3, $4, $5);

-- name: FindAPIKeyPostgresql :one
SELECT *
from api_keys
WHERE api_key = $1
LIMIT 1;

-- name: DeleteAPIKeyPostgresql :exec
DELETE FROM api_keys WHERE id = $1 AND workspace_id = $2;
