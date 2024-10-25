package admin

import (
	"github.com/gapidobri/prizer/internal/api"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	engine                     *gin.Engine
	gameService                *service.GameService
	userService                *service.UserService
	prizeService               *service.PrizeService
	wonPrizeService            *service.WonPrizeService
	participationMethodService *service.ParticipationMethodService
	drawMethodService          *service.DrawMethodService
}

func NewServer(
	db *sqlx.DB,
	gameService *service.GameService,
	userService *service.UserService,
	prizeService *service.PrizeService,
	wonPrizeService *service.WonPrizeService,
	participationMethodService *service.ParticipationMethodService,
	drawMethodService *service.DrawMethodService,
) *Server {
	return &Server{
		engine:                     api.NewServer(db),
		gameService:                gameService,
		userService:                userService,
		prizeService:               prizeService,
		wonPrizeService:            wonPrizeService,
		participationMethodService: participationMethodService,
		drawMethodService:          drawMethodService,
	}
}

func (s *Server) Run(address string) {
	s.gameRoutes()
	s.userRoutes()
	s.prizeRoutes()
	s.wonPrizeRoutes()
	s.participationMethodRoutes()
	s.drawMethodRoutes()

	log.Infof("Admin API listening on %s", address)

	err := s.engine.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
