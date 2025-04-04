package repositories

import (
	"avan/internal/db/sqlc"
	"context"
)

type WorkspaceRepository struct {
	Base Repository
}

func (w WorkspaceRepository) ListWorkspaces(ctx context.Context, ownerID int64) ([]sqlc.Workspace, error) {
	switch w.Base.GetDBDialector() {
	case "postgres":
		return w.Base.DB.Queries.ListWorkspacesPostgresql(ctx, ownerID)
	case "mysql":
		return w.Base.DB.Queries.ListWorkspacesMysql(ctx, ownerID)
	case "sqlite":
		return w.Base.DB.Queries.ListWorkspacesSqlite(ctx, ownerID)
	}
	return nil, nil
}

func (w WorkspaceRepository) SetDefaultWorkspace(ctx context.Context, id int64, ownerID int64) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.SetDefaultWorkspaceMysql(ctx, sqlc.SetDefaultWorkspaceMysqlParams{
			ID:      id,
			OwnerID: ownerID,
		})
	case "postgres":
		return w.Base.DB.Queries.SetDefaultWorkspacePostgresql(ctx, sqlc.SetDefaultWorkspacePostgresqlParams{
			ID:      int32(id),
			OwnerID: ownerID,
		})
	case "sqlite":
		return w.Base.DB.Queries.SetDefaultWorkspaceSqlite(ctx, sqlc.SetDefaultWorkspaceSqliteParams{
			ID:      id,
			OwnerID: ownerID,
		})
	}
	return nil
}

func (w WorkspaceRepository) GetDefaultWorkspace(ctx context.Context, ownerID int64) (sqlc.Workspace, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetDefaultWorkspaceMysql(ctx, ownerID)
	case "postgres":
		return w.Base.DB.Queries.GetDefaultWorkspacePostgresql(ctx, ownerID)
	case "sqlite":
		return w.Base.DB.Queries.GetDefaultWorkspaceSqlite(ctx, ownerID)
	}
	return sqlc.Workspace{}, nil
}

func (w WorkspaceRepository) UpdateWorkspace(ctx context.Context, id int64, name string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateWorkspaceMysql(ctx, sqlc.UpdateWorkspaceMysqlParams{
			ID:   id,
			Name: name,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateWorkspacePostgresql(ctx, sqlc.UpdateWorkspacePostgresqlParams{
			ID:   int32(id),
			Name: name,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateWorkspaceSqlite(ctx, sqlc.UpdateWorkspaceSqliteParams{
			ID:   id,
			Name: name,
		})
	}
	return nil
}

func (w WorkspaceRepository) DeleteWorkspace(ctx context.Context, id int64) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.DeleteWorkspaceMysql(ctx, id)
	case "postgres":
		return w.Base.DB.Queries.DeleteWorkspacePostgresql(ctx, int32(id))
	case "sqlite":
		return w.Base.DB.Queries.DeleteWorkspaceSqlite(ctx, id)
	}
	return nil
}

func (w WorkspaceRepository) CreateWorkspace(ctx context.Context, name string, ownerID int64) (bool, error) {
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		err = w.Base.DB.Queries.CreateWorkspaceMysql(ctx, sqlc.CreateWorkspaceMysqlParams{
			Name:    name,
			OwnerID: ownerID,
		})
		break
	case "postgres":
		err = w.Base.DB.Queries.CreateWorkspacePostgresql(ctx, sqlc.CreateWorkspacePostgresqlParams{
			Name:    name,
			OwnerID: ownerID,
		})
		break
	case "sqlite":
		err = w.Base.DB.Queries.CreateWorkspaceSqlite(ctx, sqlc.CreateWorkspaceSqliteParams{
			Name:    name,
			OwnerID: ownerID,
		})
		break
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (w WorkspaceRepository) FindWorkspaceByName(ctx context.Context, name string) (sqlc.Workspace, error) {
	var workspace sqlc.Workspace
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		workspace, err = w.Base.DB.Queries.FindWorkspaceByNameMysql(ctx, name)
		break
	case "postgres":
		workspace, err = w.Base.DB.Queries.FindWorkspaceByNamePostgresql(ctx, name)
		break
	case "sqlite":
		workspace, err = w.Base.DB.Queries.FindWorkspaceByNameSqlite(ctx, name)
		break
	}

	if err != nil {
		return sqlc.Workspace{}, err
	}

	return workspace, nil
}

func (w WorkspaceRepository) FindWorkspaceById(id int64) (sqlc.Workspace, error) {
	ctx := context.Background()
	var workspace sqlc.Workspace
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		workspace, err = w.Base.DB.Queries.FindWorkspaceByIdMysql(ctx, id)
		break
	case "postgres":
		workspace, err = w.Base.DB.Queries.FindWorkspaceByIdPostgresql(ctx, int32(id))
		break
	case "sqlite":
		workspace, err = w.Base.DB.Queries.FindWorkspaceByIdSqlite(ctx, id)
		break
	}

	if err != nil {
		return sqlc.Workspace{}, err
	}

	return workspace, nil
}

func (w WorkspaceRepository) FindWorkspaceByToken(ctx context.Context, token string) (sqlc.Workspace, error) {
	var workspace sqlc.Workspace
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		workspace, err = w.Base.DB.Queries.GetWorkspaceByTokenMysql(ctx, token)
		break
	case "postgres":
		workspace, err = w.Base.DB.Queries.GetWorkspaceByTokenPostgresql(ctx, token)
		break
	case "sqlite":
		workspace, err = w.Base.DB.Queries.GetWorkspaceByTokenSqlite(ctx, token)
		break
	}

	if err != nil {
		return sqlc.Workspace{}, err
	}

	return workspace, nil
}

func NewWorkspaceRepository(base Repository) WorkspaceRepository {
	return WorkspaceRepository{Base: base}
}
