-- +goose Up
CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    name         VARCHAR(50)  NOT NULL,
    email        VARCHAR(50)  NOT NULL UNIQUE,
    password     VARCHAR(255) NOT NULL,
    mobile       VARCHAR(20) UNIQUE DEFAULT NULL,
    created_at   TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP          DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workspaces
(
    id         SERIAL PRIMARY KEY,
    owner_id   BIGINT       NOT NULL,
    name       VARCHAR(255) NOT NULL UNIQUE,
    is_default BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE roles
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT      NOT NULL,
    name         VARCHAR(50) NOT NULL UNIQUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);

CREATE TABLE permissions
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    name         VARCHAR(100) NOT NULL UNIQUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);

CREATE TABLE role_permissions
(
    id            SERIAL PRIMARY KEY,
    workspace_id  BIGINT NOT NULL,
    role_id       BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE templates
(
    id           SERIAL PRIMARY KEY,
    owner_id     BIGINT      NOT NULL,
    workspace_id BIGINT      NOT NULL,
    name         VARCHAR(64) NOT NULL,
    content      TEXT        NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE messages
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    provider_id  BIGINT       NOT NULL,
    template_id  BIGINT       DEFAULT NULL,
    recipient    VARCHAR(255) NOT NULL,
    content      TEXT         NOT NULL,
    status       INT          NOT NULL,
    error        TEXT         DEFAULT NULL,
    payload      VARCHAR(255) DEFAULT NULL,
    schedule_at  TIMESTAMP    NULL,
    sent_at      TIMESTAMP    NULL,
    created_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE
);

CREATE TABLE providers
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT      NOT NULL,
    owner_id     BIGINT      NOT NULL,
    config_type  INT         NOT NULL,
    name         VARCHAR(32) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE email_configs
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    owner_id     BIGINT       NOT NULL,
    provider_id  BIGINT       NOT NULL,
    smtp_host    VARCHAR(255) NOT NULL,
    smtp_port    BIGINT       NOT NULL,
    username     VARCHAR(255) DEFAULT NULL,
    password     VARCHAR(255) DEFAULT NULL,
    from_email   VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE sms_configs
(
    id                SERIAL PRIMARY KEY,
    workspace_id      BIGINT       NOT NULL,
    owner_id          BIGINT       NOT NULL,
    provider_id       BIGINT       NOT NULL,
    url               VARCHAR(512) NOT NULL,
    method            INT          NOT NULL DEFAULT 0,
    auth_token_key    VARCHAR(64)           DEFAULT NULL,
    auth_token_value  TEXT                  DEFAULT NULL,
    is_auth_in_header BOOLEAN               DEFAULT FALSE,
    created_at        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE telegram_configs
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    owner_id     BIGINT       NOT NULL,
    provider_id  BIGINT       NOT NULL,
    bot_token    VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE slack_configs
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    owner_id     BIGINT       NOT NULL,
    provider_id  BIGINT       NOT NULL,
    webhook_url  VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE webhook_configs
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    owner_id     BIGINT       NOT NULL,
    provider_id  BIGINT       NOT NULL,
    url          VARCHAR(512) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE api_keys
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    owner_id     BIGINT       NOT NULL,
    name         VARCHAR(50)  NOT NULL,
    api_key      VARCHAR(255) NOT NULL,
    expire_at    TIMESTAMP DEFAULT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE settings
(
    id           SERIAL PRIMARY KEY,
    workspace_id BIGINT      NOT NULL,
    key          VARCHAR(32) NOT NULL,
    value        TEXT      DEFAULT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);
