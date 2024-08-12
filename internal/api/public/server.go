package public

import (
	"github.com/gapidobri/prizer/internal/api"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

func (s *Server) Run(address string) {
	s.engine.Use(api.ErrorHandler)

	s.participationMethodRoutes()

	log.Infof("Public API listening on %s", address)

	err := s.engine.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
