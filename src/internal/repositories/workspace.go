package repositories

import (
	"avan/internal/db/sqlc"
	"context"
)

type WorkspaceRepository struct {
	Base Repository
}

func (w WorkspaceRepository) ListWorkspaces(ctx context.Context, ownerID int64) ([]sqlc.Workspace, error) {
	return w.Base.DB.Queries.ListWorkspacesPostgresql(ctx, ownerID)
}

func (w WorkspaceRepository) ListAllWorkspaces(ctx context.Context) ([]sqlc.Workspace, error) {
	return w.Base.DB.Queries.ListAllWorkspacesPostgresql(ctx)
}

func (w WorkspaceRepository) SetDefaultWorkspace(ctx context.Context, id int64, ownerID int64) error {
	return w.Base.DB.Queries.SetDefaultWorkspacePostgresql(ctx, sqlc.SetDefaultWorkspacePostgresqlParams{
		ID:      id,
		OwnerID: ownerID,
	})
}

func (w WorkspaceRepository) GetDefaultWorkspace(ctx context.Context, ownerID int64) (sqlc.Workspace, error) {
	return w.Base.DB.Queries.GetDefaultWorkspacePostgresql(ctx, ownerID)
}

func (w WorkspaceRepository) UpdateWorkspace(ctx context.Context, id int64, name string) error {
	return w.Base.DB.Queries.UpdateWorkspacePostgresql(ctx, sqlc.UpdateWorkspacePostgresqlParams{
		ID:   id,
		Name: name,
	})
}

func (w WorkspaceRepository) DeleteWorkspace(ctx context.Context, id int64) error {
	return w.Base.DB.Queries.DeleteWorkspacePostgresql(ctx, id)
}

func (w WorkspaceRepository) CreateWorkspace(ctx context.Context, name string, ownerID int64) (bool, error) {
	var err error
	err = w.Base.DB.Queries.CreateWorkspacePostgresql(ctx, sqlc.CreateWorkspacePostgresqlParams{
		Name:    name,
		OwnerID: ownerID,
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (w WorkspaceRepository) FindWorkspaceByName(ctx context.Context, name string) (sqlc.Workspace, error) {
	var workspace sqlc.Workspace
	var err error
	workspace, err = w.Base.DB.Queries.FindWorkspaceByNamePostgresql(ctx, name)

	if err != nil {
		return sqlc.Workspace{}, err
	}

	return workspace, nil
}

func (w WorkspaceRepository) FindWorkspaceById(id int64) (sqlc.Workspace, error) {
	ctx := context.Background()
	var workspace sqlc.Workspace
	var err error
	workspace, err = w.Base.DB.Queries.FindWorkspaceByIdPostgresql(ctx, id)

	if err != nil {
		return sqlc.Workspace{}, err
	}

	return workspace, nil
}

func (w WorkspaceRepository) FindWorkspaceByToken(ctx context.Context, token string) (sqlc.Workspace, error) {
	var workspace sqlc.Workspace
	var err error
	workspace, err = w.Base.DB.Queries.GetWorkspaceByTokenPostgresql(ctx, token)

	if err != nil {
		return sqlc.Workspace{}, err
	}

	return workspace, nil
}

func NewWorkspaceRepository(base Repository) WorkspaceRepository {
	return WorkspaceRepository{Base: base}
}
