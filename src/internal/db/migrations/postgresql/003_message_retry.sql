-- +goose Up
ALTER TABLE messages ADD COLUMN retry BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE messages DROP COLUMN retry;
