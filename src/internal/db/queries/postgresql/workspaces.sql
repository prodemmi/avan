-- name: ListWorkspacesPostgresql :many
SELECT *
FROM workspaces
WHERE owner_id = $1 OR id IN (SELECT workspace_id FROM user_roles WHERE user_id = $1)
ORDER BY created_at DESC;

-- name: ListAllWorkspacesPostgresql :many
SELECT * FROM workspaces ORDER BY created_at DESC;

-- name: CreateWorkspacePostgresql :exec
INSERT INTO workspaces (name, owner_id) VALUES ($1, $2);

-- name: UpdateWorkspacePostgresql :exec
UPDATE workspaces SET name = $2 WHERE id = $1;

-- name: SetDefaultWorkspacePostgresql :exec
UPDATE workspaces SET is_default = (id = $1) WHERE owner_id = $2;

-- name: GetDefaultWorkspacePostgresql :one
SELECT * FROM workspaces
WHERE (owner_id = $1 AND is_default = true)
   OR id = (SELECT workspace_id FROM user_roles WHERE user_id = $1 ORDER BY created_at LIMIT 1)
ORDER BY is_default DESC
LIMIT 1;

-- name: DeleteWorkspacePostgresql :exec
DELETE FROM workspaces WHERE id = $1 AND is_default = false;

-- name: FindWorkspaceByNamePostgresql :one
SELECT * FROM workspaces WHERE name = $1 LIMIT 1;

-- name: FindWorkspaceByIdPostgresql :one
SELECT * FROM workspaces WHERE id = $1 LIMIT 1;

-- name: GetWorkspaceByTokenPostgresql :one
SELECT workspaces.id,
       workspaces.owner_id,
       workspaces.name,
       workspaces.is_default,
       workspaces.created_at,
       workspaces.updated_at,
       workspaces.business_id
FROM workspaces
         LEFT JOIN api_keys ON workspaces.id = api_keys.workspace_id
WHERE api_keys.api_key = $1
LIMIT 1;
