-- name: ListTemplatesSqlite :many
SELECT *
FROM templates
WHERE workspace_id = ? AND owner_id = ?
ORDER BY created_at DESC;

-- name: CreateTemplateSqlite :exec
INSERT INTO templates (name, content, workspace_id, owner_id)
VALUES (?, ?, ?, ?);

-- name: UpdateTemplateSqlite :exec
UPDATE templates
SET name    = ?,
    content = ?
WHERE id = ?;

-- name: DeleteTemplateSqlite :exec
DELETE
FROM templates
WHERE id = ?;

-- name: FindTemplateByNameSqlite :one
SELECT *
FROM templates
WHERE name = ?
LIMIT 1;

-- name: FindTemplateByIdSqlite :one
SELECT *
FROM templates
WHERE id = ?
LIMIT 1;