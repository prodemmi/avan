//go:build wireinject
// +build wireinject

package wire

import (
	"avan/internal/config"
	"avan/internal/db"
	"avan/internal/handlers"
	"avan/internal/logger"
	"avan/internal/network"
	"avan/internal/notifiers"
	"avan/internal/redis"
	"avan/internal/repositories"
	"avan/internal/router"
	"avan/internal/server"
	"avan/internal/services"

	"github.com/google/wire"
)

type App struct {
	Config           config.Config
	Server           server.Server
	Router           router.Router
	Notifier         notifiers.Notifier
	GossipService    network.GossipService
	UserService      services.UserService
	WorkspaceService services.WorkspaceService
	ProviderService  services.ProviderService
}

func InitializeApp(appFlags config.AppFlags) App {
	wire.Build(
		NewApp,
		router.NewRouter,
		notifiers.NewNotifier,

		handlers.NewSettingHandler,
		handlers.NewAnalysisHandler,
		handlers.NewMessageHandler,
		handlers.NewApiHandler,
		handlers.NewApiKeyHandler,
		handlers.NewUserHandler,
		handlers.NewTemplateHandler,
		handlers.NewProviderHandler,
		handlers.NewWorkspaceHandler,
		handlers.NewAuthHandler,

		services.NewSettingService,
		services.NewAnalysisService,
		services.NewApiService,
		services.NewMessageService,
		services.NewApiKeyService,
		services.NewTemplateService,
		services.NewProviderService,
		services.NewWorkspaceService,
		services.NewAuthService,
		services.NewUserService,

		repositories.NewSettingRepository,
		repositories.NewAnalysisRepository,
		repositories.NewMessageRepository,
		repositories.NewApiKeyRepository,
		repositories.NewTemplateRepository,
		repositories.NewProviderRepository,
		repositories.NewWorkspaceRepository,
		repositories.NewUserRepository,
		repositories.NewRepository,
		redis.NewRedis,
		db.NewDatabase,
		server.NewServer,
		network.NewGossipService,
		logger.NewLogger,
		config.NewConfig,
	)
	return App{}
}

func NewApp(config config.Config, server server.Server, router router.Router, notifier notifiers.Notifier, gossipService network.GossipService, userService services.UserService, workspaceService services.WorkspaceService, providerService services.ProviderService) App {
	return App{
		Config:           config,
		Server:           server,
		Router:           router,
		Notifier:         notifier,
		GossipService:    gossipService,
		UserService:      userService,
		WorkspaceService: workspaceService,
		ProviderService:  providerService,
	}
}
