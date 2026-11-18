-- name: ListMessagesSqlite :many
SELECT *
FROM messages
WHERE workspace_id = sqlc.arg('workspace_id')
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('recipient') IS NULL OR recipient = sqlc.narg('recipient'))
  AND (sqlc.narg('content') IS NULL OR content LIKE '%' || sqlc.narg('content') || '%')
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
ORDER BY created_at DESC;

-- name: CreateMessageSqlite :execlastid
INSERT INTO messages (workspace_id,
                      provider_id,
                      template_id,
                      recipient,
                      content,
                      status,
                      error,
                      schedule_at,
                      sent_at,
                      payload)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: MessageCountSqlite :one
SELECT COUNT(*)
FROM messages
WHERE workspace_id = ?
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'));

-- name: UpdateMessageStatusSqlite :exec
UPDATE messages SET status = ? WHERE id = ?;

-- name: UpdateMessageSentAtSqlite :exec
UPDATE messages SET sent_at = ? WHERE id = ?;

-- name: GetMessageSqlite :one
SELECT * FROM messages WHERE id = ? LIMIT 1;