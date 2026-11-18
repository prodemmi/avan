package repositories

import (
	"avan/internal/db/sqlc"
	"context"
)

type TemplateRepository struct {
	Base Repository
}

func (w TemplateRepository) ListTemplates(ctx context.Context, workspaceID, ownerID int64) ([]sqlc.Template, error) {
	return w.Base.DB.Queries.ListTemplatesPostgresql(ctx, sqlc.ListTemplatesPostgresqlParams{
		WorkspaceID: workspaceID,
		OwnerID:     ownerID,
	})
}

func (w TemplateRepository) UpdateTemplate(ctx context.Context, id int64, name, content string) error {
	return w.Base.DB.Queries.UpdateTemplatePostgresql(ctx, sqlc.UpdateTemplatePostgresqlParams{
		ID:      id,
		Name:    name,
		Content: content,
	})
}

func (w TemplateRepository) DeleteTemplate(ctx context.Context, id int64) error {
	return w.Base.DB.Queries.DeleteTemplatePostgresql(ctx, id)
}

func (w TemplateRepository) CreateTemplate(ctx context.Context, name string, content string, ownerID, workspaceID int64) (bool, error) {
	var err error
	err = w.Base.DB.Queries.CreateTemplatePostgresql(ctx, sqlc.CreateTemplatePostgresqlParams{
		Name:        name,
		WorkspaceID: workspaceID,
		OwnerID:     ownerID,
		Content:     content,
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (w TemplateRepository) FindTemplateByName(ctx context.Context, name string) (sqlc.Template, error) {
	var template sqlc.Template
	var err error
	template, err = w.Base.DB.Queries.FindTemplateByNamePostgresql(ctx, name)

	if err != nil {
		return sqlc.Template{}, err
	}

	return template, nil
}

func (w TemplateRepository) FindTemplateById(ctx context.Context, id int64) (sqlc.Template, error) {
	var template sqlc.Template
	var err error
	template, err = w.Base.DB.Queries.FindTemplateByIdPostgresql(ctx, id)

	if err != nil {
		return sqlc.Template{}, err
	}

	return template, nil
}

func NewTemplateRepository(base Repository) TemplateRepository {
	return TemplateRepository{Base: base}
}
