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

	switch w.Base.GetDBDialector() {
	case "postgres":
		user, err = w.Base.DB.Queries.FindUserByEmailPostgresql(ctx, email)
		break
	case "mysql":
		user, err = w.Base.DB.Queries.FindUserByEmailMysql(ctx, email)
		break
	case "sqlite":
		user, err = w.Base.DB.Queries.FindUserByEmailSqlite(ctx, email)
		break
	}

	return user, err
}

func (w UserRepository) FindUserByID(ctx context.Context, id int64) (sqlc.User, error) {
	var user sqlc.User
	var err error

	switch w.Base.GetDBDialector() {
	case "postgres":
		user, err = w.Base.DB.Queries.FindUserByIDMPostgresql(ctx, int32(id))
		break
	case "mysql":
		user, err = w.Base.DB.Queries.FindUserByIDMysql(ctx, id)
		break
	case "sqlite":
		user, err = w.Base.DB.Queries.FindUserByIDSqlite(ctx, id)
		break
	}
	return user, err
}

func (w UserRepository) CreateAdminUser(ctx context.Context, email string, password string, workspaceID int64) error {
	var err error

	switch w.Base.GetDBDialector() {
	case "postgres":
		err = w.Base.DB.Queries.CreateAdminUserPostgresql(ctx, sqlc.CreateAdminUserPostgresqlParams{
			Email:       email,
			Password:    password,
			WorkspaceID: workspaceID,
		})
		break
	case "mysql":
		err = w.Base.DB.Queries.CreateAdminUserMysql(ctx, sqlc.CreateAdminUserMysqlParams{
			Email:       email,
			Password:    password,
			WorkspaceID: workspaceID,
		})
		break
	case "sqlite":
		err = w.Base.DB.Queries.CreateAdminUserSqlite(ctx, sqlc.CreateAdminUserSqliteParams{
			Email:       email,
			Password:    password,
			WorkspaceID: workspaceID,
		})
		break
	}

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
