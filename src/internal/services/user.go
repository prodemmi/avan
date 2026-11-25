package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/repositories"
	"context"
)

func (a UserService) GetRoles(ctx context.Context, workspaceID int64) ([]sqlc.GetWorkspaceRolesPostgresqlRow, error) {
	return a.Repo.ListRoles(ctx, workspaceID)
}

func (a UserService) EnsureDefaultRoles(ctx context.Context, workspaceID int64) error {
	return a.Repo.EnsureDefaultRoles(ctx, workspaceID)
}

func (a UserService) SaveRole(ctx context.Context, workspaceID, roleID int64, name string, permissions []string) error {
	return a.Repo.SaveRole(ctx, workspaceID, roleID, name, permissions)
}

func (a UserService) DeleteRole(ctx context.Context, workspaceID, roleID int64) error {
	return a.Repo.DeleteRole(ctx, workspaceID, roleID)
}

func (a UserService) AssignRole(ctx context.Context, workspaceID, userID, roleID int64) error {
	return a.Repo.AssignRole(ctx, workspaceID, userID, roleID)
}

func (a UserService) IsSuperAdmin(ctx context.Context, userID int64) (bool, error) {
	return a.Repo.IsSuperAdmin(ctx, userID)
}

func (a UserService) UpdateUser(ctx context.Context, userID int64, name, email, password string) error {
	var hashed *string
	if password != "" {
		value := a.Repo.Base.HashPassword(password)
		hashed = &value
	}
	return a.Repo.UpdateUser(ctx, userID, name, email, hashed)
}

func (a UserService) DeleteUser(ctx context.Context, userID int64) error {
	return a.Repo.DeleteUser(ctx, userID)
}

func (a UserService) UserWorkspaceRoles(ctx context.Context, userID int64) ([]sqlc.UserRole, error) {
	return a.Repo.UserWorkspaceRoles(ctx, userID)
}

func (a UserService) AssignWorkspaces(ctx context.Context, userID int64, workspaceIDs []int64, assignments map[int64]int64) error {
	return a.Repo.AssignWorkspaces(ctx, userID, workspaceIDs, assignments)
}

func (a UserService) CreateWorkspaceUser(ctx context.Context, workspaceID int64, name, email, password string, roleID int64) error {
	return a.Repo.CreateWorkspaceUser(ctx, workspaceID, name, email, password, roleID)
}

func (a UserService) HasPermission(ctx context.Context, workspaceID, userID int64, permission string) (bool, error) {
	return a.Repo.HasPermission(ctx, workspaceID, userID, permission)
}

type UserService struct {
	Repo repositories.UserRepository
}

func (a UserService) CreateAdminUser() {
	ctx := context.Background()
	email := a.Repo.Base.Config.AdminEmail
	if email == "" {
		a.Repo.Base.Logger.Fatal("specify an admin email")
	}

	_, err := a.Repo.FindUserByEmail(ctx, email)
	if err == nil {
		return
	}

	password := a.Repo.Base.Config.AdminPassword
	if password == "" {
		a.Repo.Base.Logger.Fatal("specify an admin password")
	}

	password = a.Repo.Base.HashPassword(password)

	_ = a.Repo.CreateAdminUser(ctx, email, password, 1)
}

func (a UserService) GetDefaultUser() (sqlc.User, error) {
	ctx := context.Background()
	return a.Repo.FindUserByID(ctx, 1)
}

func (a UserService) GetUsers(ctx context.Context, workspaceID int64) ([]sqlc.GetUsersPostgresqlRow, error) {
	return a.Repo.ListUsers(ctx, workspaceID)
}

func NewUserService(repo repositories.UserRepository) UserService {
	return UserService{Repo: repo}
}
