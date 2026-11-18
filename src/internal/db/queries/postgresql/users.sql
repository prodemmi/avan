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