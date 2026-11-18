package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/repositories"
	"context"
)

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

func NewUserService(repo repositories.UserRepository) UserService {
	return UserService{Repo: repo}
}
