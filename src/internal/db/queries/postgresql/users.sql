-- name: FindUserByEmailPostgresql :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateAdminUserPostgresql :exec
INSERT INTO users (name, email, password, workspace_id)
VALUES ('admin', $1, $2, $3)
ON CONFLICT (email) DO NOTHING;

-- name: FindUserByIDMPostgresql :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUsersPostgresql :many
SELECT users.id, users.name, users.email, users.mobile, users.created_at,
       user_roles.role_id, roles.name AS role_name
FROM users
LEFT JOIN user_roles ON user_roles.user_id = users.id AND user_roles.workspace_id = $1
LEFT JOIN roles ON roles.id = user_roles.role_id AND roles.workspace_id = $1
WHERE users.id = (SELECT owner_id FROM workspaces WHERE id = $1)
   OR EXISTS (SELECT 1 FROM user_roles WHERE user_roles.workspace_id = $1 AND user_roles.user_id = users.id)
ORDER BY users.created_at DESC, users.id DESC;

-- name: GetWorkspaceRolesPostgresql :many
SELECT r.id, r.name, COALESCE(string_agg(p.name, ',' ORDER BY p.name), '') AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id AND rp.workspace_id = r.workspace_id
LEFT JOIN permissions p ON p.id = rp.permission_id AND p.workspace_id = r.workspace_id
WHERE r.workspace_id = $1
GROUP BY r.id, r.name
ORDER BY r.name;

-- name: CreateWorkspaceUserPostgresql :one
INSERT INTO users (workspace_id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id;
