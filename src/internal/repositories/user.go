package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"strings"

	"go.uber.org/zap"
)

type UserRepository struct {
	Base Repository
}

func (w UserRepository) FindUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	var user sqlc.User
	var err error

	user, err = w.Base.DB.Queries.FindUserByEmailPostgresql(ctx, email)

	return user, err
}

func (w UserRepository) FindUserByID(ctx context.Context, id int64) (sqlc.User, error) {
	var user sqlc.User
	var err error

	user, err = w.Base.DB.Queries.FindUserByIDMPostgresql(ctx, id)
	return user, err
}

func (w UserRepository) ListUsers(ctx context.Context, workspaceID int64) ([]sqlc.GetUsersPostgresqlRow, error) {
	return w.Base.DB.Queries.GetUsersPostgresql(ctx, workspaceID)
}

func (w UserRepository) ListRoles(ctx context.Context, workspaceID int64) ([]sqlc.GetWorkspaceRolesPostgresqlRow, error) {
	return w.Base.DB.Queries.GetWorkspaceRolesPostgresql(ctx, workspaceID)
}

func (w UserRepository) IsSuperAdmin(ctx context.Context, userID int64) (bool, error) {
	return w.Base.IsConfiguredAdmin(ctx, userID)
}

func (w UserRepository) UpdateUser(ctx context.Context, userID int64, name, email string, password *string) error {
	isAdmin, err := w.IsSuperAdmin(ctx, userID)
	if err != nil {
		return err
	}
	if isAdmin || (w.Base.Config.AdminEmail != "" && strings.EqualFold(email, w.Base.Config.AdminEmail)) {
		return errors.New("the configured super administrator cannot be edited")
	}
	_, err = w.Base.DB.DB.ExecContext(ctx, `UPDATE users SET name = $2, email = $3, password = COALESCE($4, password), updated_at = CURRENT_TIMESTAMP WHERE id = $1`, userID, name, email, password)
	return err
}

func (w UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	isAdmin, err := w.IsSuperAdmin(ctx, userID)
	if err != nil {
		return err
	}
	if isAdmin {
		return errors.New("the configured super administrator cannot be deleted")
	}
	var ownsWorkspace bool
	if err := w.Base.DB.DB.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM workspaces WHERE owner_id = $1)`, userID).Scan(&ownsWorkspace); err != nil {
		return err
	}
	if ownsWorkspace {
		return errors.New("transfer or delete this user's workspaces before deleting the user")
	}
	tx, err := w.Base.DB.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, table := range []string{"templates", "providers", "api_keys", "email_configs", "sms_configs", "telegram_configs", "slack_configs", "webhook_configs"} {
		query := `UPDATE ` + table + ` child SET owner_id = workspace.owner_id FROM workspaces workspace WHERE child.workspace_id = workspace.id AND child.owner_id = $1`
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			return err
		}
	}
	result, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (w UserRepository) UserWorkspaceRoles(ctx context.Context, userID int64) ([]sqlc.UserRole, error) {
	rows, err := w.Base.DB.DB.QueryContext(ctx, `SELECT workspace_id, user_id, role_id, created_at FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []sqlc.UserRole
	for rows.Next() {
		var role sqlc.UserRole
		if err := rows.Scan(&role.WorkspaceID, &role.UserID, &role.RoleID, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (w UserRepository) AssignWorkspaces(ctx context.Context, userID int64, manageableWorkspaceIDs []int64, assignments map[int64]int64) error {
	isAdmin, err := w.IsSuperAdmin(ctx, userID)
	if err != nil {
		return err
	}
	if isAdmin {
		return errors.New("the configured super administrator's access cannot be changed")
	}
	tx, err := w.Base.DB.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, workspaceID := range manageableWorkspaceIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE workspace_id = $1 AND user_id = $2`, workspaceID, userID); err != nil {
			return err
		}
		roleID, selected := assignments[workspaceID]
		if !selected {
			continue
		}
		result, err := tx.ExecContext(ctx, `INSERT INTO user_roles (workspace_id, user_id, role_id) SELECT $1, $2, id FROM roles WHERE workspace_id = $1 AND id = $3`, workspaceID, userID, roleID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("role not found in this workspace")
		}
	}
	return tx.Commit()
}

func (w UserRepository) EnsureDefaultRoles(ctx context.Context, workspaceID int64) error {
	roles := map[string][]string{
		"Administrator": {"analysis.read", "api_keys.manage", "developer.read", "messages.read", "providers.manage", "settings.manage", "templates.manage", "users.manage", "workspaces.manage"},
		"Operator":      {"analysis.read", "developer.read", "messages.read", "providers.manage", "templates.manage"},
		"Viewer":        {"analysis.read", "developer.read", "messages.read"},
	}
	tx, err := w.Base.DB.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for name, names := range roles {
		var roleID int64
		if err := tx.QueryRowContext(ctx, `INSERT INTO roles (workspace_id, name) VALUES ($1, $2) ON CONFLICT (workspace_id, name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, workspaceID, name).Scan(&roleID); err != nil {
			return err
		}
		for _, permission := range names {
			var permissionID int64
			if err := tx.QueryRowContext(ctx, `INSERT INTO permissions (workspace_id, name) VALUES ($1, $2) ON CONFLICT (workspace_id, name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, workspaceID, permission).Scan(&permissionID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions (workspace_id, role_id, permission_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, workspaceID, roleID, permissionID); err != nil {
				return err
			}
		}
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_roles (workspace_id, user_id, role_id)
		SELECT w.id, w.owner_id, r.id FROM workspaces w
		JOIN roles r ON r.workspace_id = w.id AND r.name = 'Administrator'
		WHERE w.id = $1
		ON CONFLICT (workspace_id, user_id) DO NOTHING`, workspaceID); err != nil {
		return err
	}
	return tx.Commit()
}

func (w UserRepository) SaveRole(ctx context.Context, workspaceID, roleID int64, name string, permissionNames []string) error {
	tx, err := w.Base.DB.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if roleID == 0 {
		err = tx.QueryRowContext(ctx, `INSERT INTO roles (workspace_id, name) VALUES ($1, $2) RETURNING id`, workspaceID, name).Scan(&roleID)
	} else {
		var existingID int64
		err = tx.QueryRowContext(ctx, `UPDATE roles SET name = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND workspace_id = $2 RETURNING id`, roleID, workspaceID, name).Scan(&existingID)
		roleID = existingID
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("role not found in this workspace")
		}
	}
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = $1 AND workspace_id = $2`, roleID, workspaceID); err != nil {
		return err
	}
	for _, permission := range permissionNames {
		var permissionID int64
		err = tx.QueryRowContext(ctx, `INSERT INTO permissions (workspace_id, name) VALUES ($1, $2) ON CONFLICT (workspace_id, name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, workspaceID, permission).Scan(&permissionID)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO role_permissions (workspace_id, role_id, permission_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, workspaceID, roleID, permissionID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (w UserRepository) AssignRole(ctx context.Context, workspaceID, userID, roleID int64) error {
	if roleID == 0 {
		_, err := w.Base.DB.DB.ExecContext(ctx, `DELETE FROM user_roles WHERE workspace_id = $1 AND user_id = $2`, workspaceID, userID)
		return err
	}
	result, err := w.Base.DB.DB.ExecContext(ctx, `
		INSERT INTO user_roles (workspace_id, user_id, role_id)
		SELECT $1, u.id, r.id
		FROM users u JOIN roles r ON r.workspace_id = $1
		WHERE (u.id = (SELECT owner_id FROM workspaces WHERE id = $1)
			OR EXISTS (SELECT 1 FROM user_roles existing WHERE existing.workspace_id = $1 AND existing.user_id = u.id))
			AND u.id = $2 AND r.id = $3
		ON CONFLICT (workspace_id, user_id) DO UPDATE SET role_id = EXCLUDED.role_id`, workspaceID, userID, roleID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return errors.New("user or role not found in this workspace")
	}
	return err
}

func (w UserRepository) CreateWorkspaceUser(ctx context.Context, workspaceID int64, name, email, password string, roleID int64) error {
	tx, err := w.Base.DB.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var userID int64
	if err = tx.QueryRowContext(ctx, `INSERT INTO users (workspace_id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`, workspaceID, name, email, w.Base.HashPassword(password)).Scan(&userID); err != nil {
		return err
	}
	if roleID != 0 {
		result, err := tx.ExecContext(ctx, `INSERT INTO user_roles (workspace_id, user_id, role_id) SELECT $1, $2, id FROM roles WHERE workspace_id = $1 AND id = $3`, workspaceID, userID, roleID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("role not found in this workspace")
		}
	}
	return tx.Commit()
}

func (w UserRepository) HasPermission(ctx context.Context, workspaceID, userID int64, permission string) (bool, error) {
	isAdmin, err := w.IsSuperAdmin(ctx, userID)
	if err != nil || isAdmin {
		return isAdmin, err
	}
	var allowed bool
	err = w.Base.DB.DB.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM workspaces w WHERE w.id = $1 AND w.owner_id = $2
			UNION ALL
			SELECT 1 FROM user_roles ur
			JOIN role_permissions rp ON rp.workspace_id = ur.workspace_id AND rp.role_id = ur.role_id
			JOIN permissions p ON p.workspace_id = rp.workspace_id AND p.id = rp.permission_id
			WHERE ur.workspace_id = $1 AND ur.user_id = $2 AND p.name = $3
		)`, workspaceID, userID, permission).Scan(&allowed)
	return allowed, err
}

func (w UserRepository) CreateAdminUser(ctx context.Context, email string, password string, workspaceID int64) error {
	var err error

	err = w.Base.DB.Queries.CreateAdminUserPostgresql(ctx, sqlc.CreateAdminUserPostgresqlParams{
		Email:       email,
		Password:    password,
		WorkspaceID: workspaceID,
	})

	if err != nil {
		w.Base.Logger.Error("failed to create admin user",
			zap.String("event", "create_admin_user"),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func NewUserRepository(base Repository) UserRepository {
	return UserRepository{Base: base}
}
