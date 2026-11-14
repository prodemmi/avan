package server

import (
	"avan/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Server struct {
	Config config.Config
	Engine *gin.Engine
	Logger zap.Logger
}

func (s Server) Start() {
	port := strconv.Itoa(s.Config.ServerPort)
	if port == "" {
		port = "8080"
	}

	if err := s.Engine.Run(":" + port); err != nil {
		s.Logger.Fatal("failed to start server", zap.Error(err))
	}
}

func NewServer(config config.Config, logger zap.Logger) Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	return Server{Config: config, Engine: engine, Logger: logger}
}
