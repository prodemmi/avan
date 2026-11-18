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
SELECT id, name, email, mobile, created_at
FROM users
WHERE workspace_id = $1
ORDER BY created_at DESC, id DESC;
