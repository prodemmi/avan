-- name: ListTemplatesMysql :many
SELECT *
FROM templates
WHERE workspace_id = ? AND owner_id = ?
ORDER BY created_at DESC;

-- name: CreateTemplateMysql :exec
INSERT INTO templates (name, content, workspace_id, owner_id)
VALUES (?, ?, ?, ?);

-- name: UpdateTemplateMysql :exec
UPDATE templates
SET name    = ?,
    content = ?
WHERE id = ?;

-- name: DeleteTemplateMysql :exec
DELETE
FROM templates
WHERE id = ?;

-- name: FindTemplateByNameMysql :one
SELECT *
FROM templates
WHERE name = ?
LIMIT 1;

-- name: FindTemplateByIdMysql :one
SELECT *
FROM templates
WHERE id = ?
LIMIT 1;