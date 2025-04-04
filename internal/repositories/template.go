package repositories

import (
	"avan/internal/db/sqlc"
	"context"
)

type TemplateRepository struct {
	Base Repository
}

func (w TemplateRepository) ListTemplates(ctx context.Context, workspaceID, ownerID int64) ([]sqlc.Template, error) {
	switch w.Base.GetDBDialector() {
	case "postgres":
		return w.Base.DB.Queries.ListTemplatesPostgresql(ctx, sqlc.ListTemplatesPostgresqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	case "mysql":
		return w.Base.DB.Queries.ListTemplatesMysql(ctx, sqlc.ListTemplatesMysqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	case "sqlite":
		return w.Base.DB.Queries.ListTemplatesSqlite(ctx, sqlc.ListTemplatesSqliteParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	}
	return nil, nil
}

func (w TemplateRepository) UpdateTemplate(ctx context.Context, id int64, name, content string) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateTemplateMysql(ctx, sqlc.UpdateTemplateMysqlParams{
			ID:      id,
			Name:    name,
			Content: content,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateTemplatePostgresql(ctx, sqlc.UpdateTemplatePostgresqlParams{
			ID:      int32(id),
			Name:    name,
			Content: content,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateTemplateSqlite(ctx, sqlc.UpdateTemplateSqliteParams{
			ID:      id,
			Name:    name,
			Content: content,
		})
	}
	return nil
}

func (w TemplateRepository) DeleteTemplate(ctx context.Context, id int64) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.DeleteTemplateMysql(ctx, id)
	case "postgres":
		return w.Base.DB.Queries.DeleteTemplatePostgresql(ctx, int32(id))
	case "sqlite":
		return w.Base.DB.Queries.DeleteTemplateSqlite(ctx, id)
	}
	return nil
}

func (w TemplateRepository) CreateTemplate(ctx context.Context, name string, content string, ownerID, workspaceID int64) (bool, error) {
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		err = w.Base.DB.Queries.CreateTemplateMysql(ctx, sqlc.CreateTemplateMysqlParams{
			Name:        name,
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			Content:     content,
		})
		break
	case "postgres":
		err = w.Base.DB.Queries.CreateTemplatePostgresql(ctx, sqlc.CreateTemplatePostgresqlParams{
			Name:        name,
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			Content:     content,
		})
		break
	case "sqlite":
		err = w.Base.DB.Queries.CreateTemplateSqlite(ctx, sqlc.CreateTemplateSqliteParams{
			Name:        name,
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			Content:     content,
		})
		break
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (w TemplateRepository) FindTemplateByName(ctx context.Context, name string) (sqlc.Template, error) {
	var template sqlc.Template
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		template, err = w.Base.DB.Queries.FindTemplateByNameMysql(ctx, name)
		break
	case "postgres":
		template, err = w.Base.DB.Queries.FindTemplateByNamePostgresql(ctx, name)
		break
	case "sqlite":
		template, err = w.Base.DB.Queries.FindTemplateByNameSqlite(ctx, name)
		break
	}

	if err != nil {
		return sqlc.Template{}, err
	}

	return template, nil
}

func (w TemplateRepository) FindTemplateById(ctx context.Context, id int64) (sqlc.Template, error) {
	var template sqlc.Template
	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		template, err = w.Base.DB.Queries.FindTemplateByIdMysql(ctx, id)
		break
	case "postgres":
		template, err = w.Base.DB.Queries.FindTemplateByIdPostgresql(ctx, int32(id))
		break
	case "sqlite":
		template, err = w.Base.DB.Queries.FindTemplateByIdSqlite(ctx, id)
		break
	}

	if err != nil {
		return sqlc.Template{}, err
	}

	return template, nil
}

func NewTemplateRepository(base Repository) TemplateRepository {
	return TemplateRepository{Base: base}
}
