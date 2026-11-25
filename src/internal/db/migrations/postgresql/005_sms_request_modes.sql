-- +goose Up
ALTER TABLE sms_configs
    ADD COLUMN body_encoding VARCHAR(16) NOT NULL DEFAULT 'json'
        CHECK (body_encoding IN ('json', 'form')),
    ADD COLUMN auth_type VARCHAR(16) NOT NULL DEFAULT 'token'
        CHECK (auth_type IN ('token', 'basic')),
    ADD COLUMN auth_username VARCHAR(255),
    ADD COLUMN account_sid VARCHAR(64),
    ADD COLUMN sender VARCHAR(255);

-- +goose Down
ALTER TABLE sms_configs
    DROP COLUMN body_encoding,
    DROP COLUMN auth_type,
    DROP COLUMN auth_username,
    DROP COLUMN account_sid,
    DROP COLUMN sender;
