-- name: ListTemplatesPostgresql :many
SELECT *
FROM templates
WHERE workspace_id = $1 AND owner_id = $2
ORDER BY created_at DESC;

-- name: CreateTemplatePostgresql :exec
INSERT INTO templates (name, content, workspace_id, owner_id)
VALUES ($1, $2, $3, $4);

-- name: UpdateTemplatePostgresql :exec
UPDATE templates
SET name = $2, content = $3
WHERE id = $1;

-- name: DeleteTemplatePostgresql :exec
DELETE
FROM templates
WHERE id = $1;

-- name: FindTemplateByNamePostgresql :one
SELECT *
FROM templates
WHERE name = $1
LIMIT 1;

-- name: FindTemplateByIdPostgresql :one
SELECT *
FROM templates
WHERE id = $1
LIMIT 1;