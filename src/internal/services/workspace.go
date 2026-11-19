package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/redis"
	"avan/internal/repositories"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type WorkspaceService struct {
	Repo  repositories.WorkspaceRepository
	Cache redis.Cache[sqlc.Workspace]

	SettingService SettingService
}

func (w WorkspaceService) GetWorkspaces(ctx context.Context, ownerID int64) ([]sqlc.Workspace, error) {
	isAdmin, err := w.Repo.Base.IsConfiguredAdmin(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		return w.Repo.ListAllWorkspaces(ctx)
	}
	return w.Repo.ListWorkspaces(ctx, ownerID)
}

func (w WorkspaceService) CreateWorkspace(name string, ownerID int64) (bool, error) {
	ctx := context.Background()

	result, err := w.Repo.CreateWorkspace(ctx, name, ownerID)
	if err == nil {
		key := w.Cache.GetStringKey("workspaces", ownerID)
		w.Cache.Delete(key...)
	}

	workspace, err := w.GetWorkspaceByName(ctx, name)
	fmt.Println(workspace, err)
	if err != nil {
		return false, err
	}

	w.SettingService.CreateDefaultSettings(workspace.ID)
	if err := w.SettingService.UserService.EnsureDefaultRoles(ctx, workspace.ID); err != nil {
		return false, err
	}

	return result, err
}

func (w WorkspaceService) SetDefaultWorkspace(id int64, ownerID int64) error {
	ctx := context.Background()
	err := w.Repo.SetDefaultWorkspace(ctx, id, ownerID)
	if err == nil {
		key := w.Cache.GetStringKey("workspace", "default", ownerID)
		w.Cache.Delete(key...)
	}

	return err
}

func (w WorkspaceService) UpdateWorkspace(id int64, name string) error {
	ctx := context.Background()
	err := w.Repo.UpdateWorkspace(ctx, id, name)
	if err == nil {
		key := w.Cache.GetStringKey("workspace", id)
		w.Cache.Delete(key...)
	}
	return err
}

func (w WorkspaceService) DeleteWorkspace(id int64) error {
	ctx := context.Background()
	err := w.Repo.DeleteWorkspace(ctx, id)
	if err == nil {
		key := w.Cache.GetStringKey("workspace", id)
		w.Cache.Delete(key...)
	}
	return err
}

func (w WorkspaceService) GetDefaultWorkspace(ownerID int64) (sqlc.Workspace, error) {
	key := w.Cache.GetKey("workspace", "default", ownerID)
	result, err := w.Cache.GetOrSet(key, time.Hour*24, func() (sqlc.Workspace, error) {
		ctx := context.Background()
		return w.Repo.GetDefaultWorkspace(ctx, ownerID)
	})
	return result, err
}

func (w WorkspaceService) GetWorkspaceByName(ctx context.Context, name string) (sqlc.Workspace, error) {
	return w.Repo.FindWorkspaceByName(ctx, name)
}

func (w WorkspaceService) GetWorkspace(id int64) (sqlc.Workspace, error) {
	key := w.Cache.GetKey("workspace", id)
	result, err := w.Cache.GetOrSet(key, time.Hour*24, func() (sqlc.Workspace, error) {
		return w.Repo.FindWorkspaceById(id)
	})
	return result, err
}

func (w WorkspaceService) CreateDefaultWorkspace() {
	ctx := context.Background()

	logError := func(err error) {
		w.Repo.Base.Logger.Fatal("failed to create default workspace", zap.Error(err))
	}

	defaultName := w.Repo.Base.Config.WorkspaceDefaultName
	_, err := w.Repo.FindWorkspaceByName(ctx, defaultName)
	if err == nil {
		return
	}

	_, err = w.CreateWorkspace(defaultName, 1)
	if err != nil {
		logError(err)
		return
	}

	err = w.SetDefaultWorkspace(1, 1)
	if err != nil {
		logError(err)
		return
	}
}

func (w WorkspaceService) GetWorkspaceByToken(token string) (sqlc.Workspace, error) {
	hashedToken := redis.HashKey(token)
	key := w.Cache.GetKey("workspace", "by_token", hashedToken)

	result, err := w.Cache.GetOrSet(key, time.Hour*8, func() (sqlc.Workspace, error) {
		ctx := context.Background()
		return w.Repo.FindWorkspaceByToken(ctx, token)
	})
	return result, err
}

func (w WorkspaceService) GetPrimaryColor(ownerID int64) (string, error) {
	workspace, err := w.GetDefaultWorkspace(ownerID)
	if err != nil {
		return "", err
	}

	setting, err := w.SettingService.GetSetting(workspace.ID, "primary-color")
	if err != nil {
		return "", err
	}

	if setting.Value != nil {
		return *setting.Value, nil
	}
	return "", nil
}

func (w WorkspaceService) GetTheme(ownerID int64) (string, error) {
	workspace, err := w.GetDefaultWorkspace(ownerID)
	if err != nil {
		return "", err
	}

	setting, err := w.SettingService.GetSetting(workspace.ID, "theme")
	if err != nil {
		return "", err
	}

	if setting.Value != nil {
		return *setting.Value, nil
	}
	return "", nil
}

func (w WorkspaceService) GetTimezone(ownerID int64) (string, error) {
	workspace, err := w.GetDefaultWorkspace(ownerID)
	if err != nil {
		return "", err
	}

	setting, err := w.SettingService.GetSetting(workspace.ID, "timezone")
	if err != nil {
		return "", err
	}

	if setting.Value != nil {
		return *setting.Value, nil
	}
	return "", nil
}

func NewWorkspaceService(repo repositories.WorkspaceRepository, settingService SettingService) WorkspaceService {
	cache := redis.NewCache[sqlc.Workspace](repo.Base.Redis)
	return WorkspaceService{Repo: repo, Cache: cache, SettingService: settingService}
}
