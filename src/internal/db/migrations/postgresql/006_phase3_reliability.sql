-- +goose Up
ALTER TABLE messages ADD COLUMN idempotency_key TEXT;
ALTER TABLE messages ADD COLUMN request_hash TEXT;
CREATE UNIQUE INDEX messages_workspace_id_idempotency_key_uq
    ON messages (workspace_id, idempotency_key);

INSERT INTO settings (workspace_id, key, value)
SELECT workspaces.id, 'api-security-policy', '{}'
FROM workspaces
WHERE NOT EXISTS (
    SELECT 1
    FROM settings
    WHERE settings.workspace_id = workspaces.id
      AND settings.key = 'api-security-policy'
);

-- +goose Down
DROP INDEX messages_workspace_id_idempotency_key_uq;
ALTER TABLE messages DROP COLUMN request_hash;
ALTER TABLE messages DROP COLUMN idempotency_key;
