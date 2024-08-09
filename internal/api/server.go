package api

import (
	"github.com/gapidobri/prizer/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine      *gin.Engine
	gameService *service.GameService
}

func NewServer(gameService *service.GameService) *Server {
	return &Server{
		engine:      gin.Default(),
		gameService: gameService,
	}
}

func (s *Server) Run() error {
	s.engine.Use(ErrorHandler)

	s.gameRoutes()

	return s.engine.Run()
}
