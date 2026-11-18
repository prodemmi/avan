package repositories

import (
	"avan/internal/db/sqlc"
	"context"

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
