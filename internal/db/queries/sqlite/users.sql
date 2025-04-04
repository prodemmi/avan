-- name: FindUserByEmailSqlite :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- name: CreateAdminUserSqlite :exec
INSERT INTO users (name, email, password, workspace_id)
VALUES ('admin', ?, ?, ?)
ON CONFLICT(email) DO NOTHING;

-- name: FindUserByIDSqlite :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;