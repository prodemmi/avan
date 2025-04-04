-- name: FindUserByEmailMysql :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- name: CreateAdminUserMysql :exec
INSERT INTO users (name, email, password, workspace_id)
VALUES ('admin', ?, ?, ?)
ON DUPLICATE KEY UPDATE id = id;

-- name: FindUserByIDMysql :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;