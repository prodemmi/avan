-- name: ListMessagesPostgresql :many
SELECT *
FROM messages
WHERE workspace_id = sqlc.arg('workspace_id')
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('recipient') IS NULL OR recipient = sqlc.narg('recipient'))
  AND (sqlc.narg('content') IS NULL OR content ILIKE '%' || sqlc.narg('content') || '%')
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
ORDER BY created_at DESC;

-- name: CreateMessagePostgresql :one
INSERT INTO messages (workspace_id,
                      provider_id,
                      template_id,
                      recipient,
                      content,
                      status,
                      error,
                      retry,
                      schedule_at,
                      sent_at,
                      payload)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id::BIGINT;

-- name: MessageCountPostgresql :one
SELECT COUNT(*)
FROM messages
WHERE workspace_id = $1
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'));

-- name: UpdateMessageStatusPostgresql :exec
UPDATE messages SET status = $2 WHERE id = $1;

-- name: UpdateMessageSentAtPostgresql :exec
UPDATE messages SET sent_at = $2 WHERE id = $1;

-- name: GetMessagePostgresql :one
SELECT * FROM messages WHERE id = $1 LIMIT 1;
