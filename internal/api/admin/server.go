package admin

import (
	"github.com/gapidobri/prizer/internal/api"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	engine          *gin.Engine
	gameService     *service.GameService
	userService     *service.UserService
	prizeService    *service.PrizeService
	wonPrizeService *service.WonPrizeService
}

func NewServer(
	gameService *service.GameService,
	userService *service.UserService,
	prizeService *service.PrizeService,
	wonPrizeService *service.WonPrizeService,
) *Server {
	return &Server{
		engine:          gin.Default(),
		gameService:     gameService,
		userService:     userService,
		prizeService:    prizeService,
		wonPrizeService: wonPrizeService,
	}
}

func (s *Server) Run(address string) {
	s.engine.Use(api.ErrorHandler)

	s.gameRoutes()
	s.userRoutes()
	s.prizeRoutes()
	s.wonPrizeRoutes()

	err := s.engine.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
