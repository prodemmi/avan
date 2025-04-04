-- name: ListMessagesMysql :many
SELECT *
FROM messages
WHERE workspace_id = ?
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT 10 OFFSET ?;

-- name: CreateMessageMysql :execlastid
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

-- name: MessageCountMysql :one
SELECT COUNT(*)
FROM messages
WHERE workspace_id = ?
  AND (sqlc.narg('template_id') IS NULL OR template_id = sqlc.narg('template_id'))
  AND (sqlc.narg('provider_id') IS NULL OR provider_id = sqlc.narg('provider_id'))
  AND (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'));

-- name: UpdateMessageStatusMysql :exec
UPDATE messages SET status = ? WHERE id = ?;

-- name: UpdateMessageSentAtMysql :exec
UPDATE messages SET sent_at = ? WHERE id = ?;

-- name: GetMessageMysql :one
SELECT * FROM messages WHERE id = ? LIMIT 1;