-- +goose Up
ALTER TABLE providers
    ADD COLUMN security_policy TEXT NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE providers DROP COLUMN security_policy;
