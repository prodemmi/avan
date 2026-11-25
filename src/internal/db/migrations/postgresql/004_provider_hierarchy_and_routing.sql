-- +goose Up
CREATE TABLE businesses
(
    id         BIGSERIAL PRIMARY KEY,
    owner_id   BIGINT       NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (id, owner_id),
    UNIQUE (owner_id, name)
);

ALTER TABLE workspaces ADD COLUMN business_id BIGINT;
INSERT INTO businesses (owner_id, name)
SELECT owner_id, name FROM workspaces;
UPDATE workspaces AS w
SET business_id = b.id
FROM businesses AS b
WHERE b.owner_id = w.owner_id AND b.name = w.name;
ALTER TABLE workspaces ALTER COLUMN business_id SET NOT NULL;
ALTER TABLE workspaces
    ADD CONSTRAINT workspaces_business_owner_fkey
    FOREIGN KEY (business_id, owner_id) REFERENCES businesses (id, owner_id) ON DELETE CASCADE;

CREATE TABLE apps
(
    id           BIGSERIAL PRIMARY KEY,
    workspace_id BIGINT       NOT NULL,
    name         VARCHAR(64)  NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    UNIQUE (workspace_id, name),
    UNIQUE (id, workspace_id)
);

INSERT INTO apps (workspace_id, name)
SELECT id, 'Default' FROM workspaces;

ALTER TABLE providers ADD COLUMN app_id BIGINT;
UPDATE providers AS p
SET app_id = a.id
FROM apps AS a
WHERE a.workspace_id = p.workspace_id AND a.name = 'Default';
ALTER TABLE providers ALTER COLUMN app_id SET NOT NULL;
ALTER TABLE providers
    ADD CONSTRAINT providers_app_workspace_fkey
    FOREIGN KEY (app_id, workspace_id) REFERENCES apps (id, workspace_id) ON DELETE CASCADE;
ALTER TABLE providers
    ADD COLUMN priority INTEGER NOT NULL DEFAULT 0 CHECK (priority >= 0),
    ADD COLUMN weight INTEGER NOT NULL DEFAULT 1 CHECK (weight BETWEEN 1 AND 1000);

CREATE FUNCTION create_workspace_business() RETURNS trigger AS $$
BEGIN
    IF NEW.business_id IS NULL THEN
        INSERT INTO businesses (owner_id, name)
        VALUES (NEW.owner_id, NEW.name)
        RETURNING id INTO NEW.business_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER workspaces_create_business
BEFORE INSERT ON workspaces
FOR EACH ROW EXECUTE FUNCTION create_workspace_business();

CREATE FUNCTION create_workspace_default_app() RETURNS trigger AS $$
BEGIN
    INSERT INTO apps (workspace_id, name) VALUES (NEW.id, 'Default');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER workspaces_create_default_app
AFTER INSERT ON workspaces
FOR EACH ROW EXECUTE FUNCTION create_workspace_default_app();

-- +goose Down
DROP TRIGGER workspaces_create_default_app ON workspaces;
DROP FUNCTION create_workspace_default_app();
DROP TRIGGER workspaces_create_business ON workspaces;
DROP FUNCTION create_workspace_business();

ALTER TABLE providers DROP CONSTRAINT providers_app_workspace_fkey;
ALTER TABLE providers DROP COLUMN app_id;
ALTER TABLE providers DROP COLUMN priority, DROP COLUMN weight;
DROP TABLE apps;

ALTER TABLE workspaces DROP CONSTRAINT workspaces_business_owner_fkey;
ALTER TABLE workspaces DROP COLUMN business_id;
DROP TABLE businesses;
