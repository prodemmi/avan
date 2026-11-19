-- +goose Up
ALTER TABLE sms_configs ADD COLUMN body_schema TEXT;

ALTER TABLE roles DROP CONSTRAINT roles_name_key;
ALTER TABLE roles ADD CONSTRAINT roles_workspace_name_key UNIQUE (workspace_id, name);
ALTER TABLE permissions DROP CONSTRAINT permissions_name_key;
ALTER TABLE permissions ADD CONSTRAINT permissions_workspace_name_key UNIQUE (workspace_id, name);

CREATE TABLE user_roles (
    workspace_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, user_id),
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE
);

ALTER TABLE role_permissions
    ADD CONSTRAINT role_permissions_workspace_role_permission_key UNIQUE (workspace_id, role_id, permission_id);

INSERT INTO roles (workspace_id, name)
SELECT workspace.id, default_role.name
FROM workspaces workspace
CROSS JOIN (VALUES ('Administrator'), ('Operator'), ('Viewer')) AS default_role(name)
ON CONFLICT (workspace_id, name) DO NOTHING;

INSERT INTO permissions (workspace_id, name)
SELECT workspace.id, permission.name
FROM workspaces workspace
CROSS JOIN (VALUES
    ('analysis.read'), ('api_keys.manage'), ('developer.read'), ('messages.read'),
    ('providers.manage'), ('settings.manage'), ('templates.manage'),
    ('users.manage'), ('workspaces.manage')
) AS permission(name)
ON CONFLICT (workspace_id, name) DO NOTHING;

INSERT INTO role_permissions (workspace_id, role_id, permission_id)
SELECT r.workspace_id, r.id, p.id
FROM roles r
JOIN permissions p ON p.workspace_id = r.workspace_id
WHERE r.name = 'Administrator'
   OR (r.name = 'Operator' AND p.name IN ('analysis.read', 'developer.read', 'messages.read', 'providers.manage', 'templates.manage'))
   OR (r.name = 'Viewer' AND p.name IN ('analysis.read', 'developer.read', 'messages.read'))
ON CONFLICT (workspace_id, role_id, permission_id) DO NOTHING;

INSERT INTO user_roles (workspace_id, user_id, role_id)
SELECT users.workspace_id, users.id, roles.id
FROM users
JOIN roles ON roles.workspace_id = users.workspace_id AND roles.name = 'Administrator'
ON CONFLICT (workspace_id, user_id) DO NOTHING;

-- +goose Down
DROP TABLE user_roles;
ALTER TABLE role_permissions DROP CONSTRAINT role_permissions_workspace_role_permission_key;
ALTER TABLE permissions DROP CONSTRAINT permissions_workspace_name_key;
ALTER TABLE permissions ADD CONSTRAINT permissions_name_key UNIQUE (name);
ALTER TABLE roles DROP CONSTRAINT roles_workspace_name_key;
ALTER TABLE roles ADD CONSTRAINT roles_name_key UNIQUE (name);
ALTER TABLE sms_configs DROP COLUMN body_schema;
