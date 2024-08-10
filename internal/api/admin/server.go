package admin

import (
	"github.com/gapidobri/prizer/internal/api"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	engine       *gin.Engine
	gameService  *service.GameService
	userService  *service.UserService
	prizeService *service.PrizeService
}

func NewServer(
	gameService *service.GameService,
	userService *service.UserService,
	prizeService *service.PrizeService,
) *Server {
	return &Server{
		engine:       gin.Default(),
		gameService:  gameService,
		userService:  userService,
		prizeService: prizeService,
	}
}

func (s *Server) Run(address string) {
	s.engine.Use(api.ErrorHandler)

	s.gameRoutes()
	s.userRoutes()
	s.prizeRoutes()

	err := s.engine.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
