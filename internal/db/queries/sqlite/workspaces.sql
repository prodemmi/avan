-- name: ListWorkspacesSqlite :many
SELECT *
FROM workspaces
WHERE owner_id = ?
ORDER BY created_at DESC;

-- name: CreateWorkspaceSqlite :exec
INSERT INTO workspaces (name, owner_id) VALUES (?, ?);

-- name: UpdateWorkspaceSqlite :exec
UPDATE workspaces SET name = ? WHERE id = ?;

-- name: SetDefaultWorkspaceSqlite :exec
UPDATE workspaces SET is_default = (id = ?) WHERE owner_id = ?;

-- name: GetDefaultWorkspaceSqlite :one
SELECT * FROM workspaces WHERE is_default = true AND owner_id = ? LIMIT 1;

-- name: DeleteWorkspaceSqlite :exec
DELETE FROM workspaces WHERE id = ? AND is_default = 0;

-- name: FindWorkspaceByNameSqlite :one
SELECT * FROM workspaces WHERE name = ? LIMIT 1;

-- name: FindWorkspaceByIdSqlite :one
SELECT * FROM workspaces WHERE id = ? LIMIT 1;

-- name: GetWorkspaceByTokenSqlite :one
SELECT workspaces.id,
       workspaces.owner_id,
       workspaces.name,
       workspaces.is_default,
       workspaces.created_at,
       workspaces.updated_at
FROM workspaces
         LEFT JOIN api_keys ON workspaces.id = api_keys.workspace_id
WHERE api_keys.api_key = ?
LIMIT 1;
