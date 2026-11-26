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

-- name: CreateIdempotentMessagePostgresql :one
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
                      payload,
                      idempotency_key,
                      request_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (workspace_id, idempotency_key) DO NOTHING
RETURNING *;

-- name: GetMessageByIdempotencyKeyPostgresql :one
SELECT * FROM messages
WHERE workspace_id = $1 AND idempotency_key = $2
LIMIT 1;

-- name: GetMessageStatusForWorkspacePostgresql :one
SELECT id, provider_id, status, schedule_at, sent_at, created_at, updated_at FROM messages
WHERE id = $1 AND workspace_id = $2
LIMIT 1;

-- name: UpdateMessageStatusForWorkspacePostgresql :execrows
UPDATE messages
SET status = $3,
    sent_at = CASE WHEN $3 IN (0, 3, 4) AND sent_at IS NULL THEN CURRENT_TIMESTAMP ELSE sent_at END,
    schedule_at = CASE WHEN $3 IN (0, 3, 4) THEN NULL ELSE schedule_at END,
    error = CASE WHEN $3 IN (0, 3, 4) THEN NULL ELSE error END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND workspace_id = $2
  AND (status = $3
       OR ($3 = 0 AND status = 1 AND schedule_at IS NULL)
       OR ($3 = 2 AND (status = 0 OR (status = 1 AND schedule_at IS NULL)))
       OR ($3 = 3 AND (status IN (0, 2) OR (status = 1 AND schedule_at IS NULL)))
       OR ($3 = 4 AND status = 3));

-- name: MessageCountPostgresql :one
SELECT COUNT(*)
FROM messages
WHERE workspace_id = $1
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'));

-- name: UpdateMessageStatusPostgresql :exec
UPDATE messages
SET status = $2
WHERE id = $1
  AND (status = $2
       OR ($2 = 0 AND status = 1 AND schedule_at IS NULL)
       OR ($2 = 1 AND status = 1)
       OR ($2 = 2 AND status = 1 AND schedule_at IS NULL));

-- name: UpdateMessageSentAtPostgresql :exec
UPDATE messages SET sent_at = $2 WHERE id = $1;

-- name: GetMessagePostgresql :one
SELECT * FROM messages WHERE id = $1 LIMIT 1;
