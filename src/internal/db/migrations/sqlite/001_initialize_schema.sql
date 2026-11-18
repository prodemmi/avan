-- +goose Up

CREATE TABLE users
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    name         TEXT NOT NULL,
    email        TEXT NOT NULL UNIQUE,
    password     TEXT NOT NULL,
    mobile       TEXT      DEFAULT NULL UNIQUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workspaces
(
    id         INT PRIMARY KEY AUTOINCREMENT,
    owner_id   INT  NOT NULL,
    name       TEXT NOT NULL UNIQUE,
    is_default INT  NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE roles
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    name         TEXT NOT NULL UNIQUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);

CREATE TABLE permissions
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    name         TEXT NOT NULL UNIQUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);

CREATE TABLE role_permissions
(
    id            INT PRIMARY KEY AUTOINCREMENT,
    workspace_id  INT NOT NULL,
    role_id       INT NOT NULL,
    permission_id INT NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id),
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE templates
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    owner_id     INT  NOT NULL,
    name         TEXT NOT NULL,
    content      TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE messages
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT       NOT NULL,
    provider_id  INT       NOT NULL,
    template_id  INT          DEFAULT NULL,
    recipient    TEXT      NOT NULL,
    content      TEXT      NOT NULL,
    status       INT       NOT NULL,
    error        TEXT         DEFAULT NULL,
    payload      VARCHAR(255) DEFAULT NULL,
    schedule_at  TIMESTAMP NULL,
    sent_at      TIMESTAMP NULL,
    created_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE
);

CREATE TABLE providers
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    owner_id     INT  NOT NULL,
    config_type  INT  NOT NULL,
    name         TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE email_configs
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    owner_id     INT  NOT NULL,
    provider_id  INT  NOT NULL,
    smtp_host    TEXT NOT NULL,
    smtp_port    INT  NOT NULL,
    username     TEXT      DEFAULT NULL,
    password     TEXT      DEFAULT NULL,
    from_email   TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE sms_configs
(
    id                INT PRIMARY KEY AUTOINCREMENT,
    workspace_id      INT          NOT NULL,
    owner_id          INT          NOT NULL,
    provider_id       INT          NOT NULL,
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
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    owner_id     INT  NOT NULL,
    provider_id  INT  NOT NULL,
    bot_token    TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE slack_configs
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT  NOT NULL,
    owner_id     INT  NOT NULL,
    provider_id  INT  NOT NULL,
    webhook_url  TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE webhook_configs
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT          NOT NULL,
    owner_id     INT          NOT NULL,
    provider_id  INT          NOT NULL,
    url          VARCHAR(512) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE api_keys
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT         NOT NULL,
    owner_id     INT         NOT NULL,
    name         VARCHAR(50) NOT NULL,
    api_key      TEXT        NOT NULL,
    expire_at    TIMESTAMP DEFAULT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE settings
(
    id           INT PRIMARY KEY AUTOINCREMENT,
    workspace_id INT NOT NULL,
    key          VARCHAR(32) NOT NULL,
    value        TEXT      DEFAULT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);
