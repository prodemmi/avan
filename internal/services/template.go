package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/redis"
	"avan/internal/repositories"
	"context"
	"time"
)

type TemplateService struct {
	Repo  repositories.TemplateRepository
	Cache redis.Cache[sqlc.Template]
}

func (w TemplateService) GetTemplates(workspaceID, ownerID int64) ([]sqlc.Template, error) {
	key := w.Cache.GetKey("templates", ownerID)
	return w.Cache.GetOrSetArray(key, time.Hour*24, func() ([]sqlc.Template, error) {
		ctx := context.Background()
		return w.Repo.ListTemplates(ctx, workspaceID, ownerID)
	})
}

func (w TemplateService) CreateTemplate(name, content string, ownerID, workspaceID int64) (bool, error) {
	ctx := context.Background()
	result, err := w.Repo.CreateTemplate(ctx, name, content, ownerID, workspaceID)

	if err == nil {
		key := w.Cache.GetStringKey("template", ownerID)
		w.Cache.Delete(key...)
	}

	return result, err
}

func (w TemplateService) UpdateTemplate(id int64, name, content string) error {
	ctx := context.Background()
	err := w.Repo.UpdateTemplate(ctx, id, name, content)
	if err == nil {
		key := w.Cache.GetStringKey("template", id)
		w.Cache.Delete(key...)
	}

	return err
}

func (w TemplateService) DeleteTemplate(id int64) error {
	ctx := context.Background()
	err := w.Repo.DeleteTemplate(ctx, id)
	if err == nil {
		key := w.Cache.GetStringKey("template", id)
		w.Cache.Delete(key...)
	}

	return err
}

func (w TemplateService) GetTemplateByName(name string) (sqlc.Template, error) {
	ctx := context.Background()
	return w.Repo.FindTemplateByName(ctx, name)
}

func (w TemplateService) GetTemplate(id int64) (sqlc.Template, error) {
	key := w.Cache.GetKey("template", id)
	return w.Cache.GetOrSet(key, time.Hour*24, func() (sqlc.Template, error) {
		ctx := context.Background()
		return w.Repo.FindTemplateById(ctx, id)
	})
}

func NewTemplateService(repo repositories.TemplateRepository) TemplateService {
	cache := redis.NewCache[sqlc.Template](repo.Base.Redis)
	return TemplateService{Repo: repo, Cache: cache}
}
