package internal

import (
	"context"
	"github.com/gapidobri/prizer/internal/api/admin"
	"github.com/gapidobri/prizer/internal/api/public"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/clients/addressvalidation"
	"github.com/gapidobri/prizer/internal/pkg/clients/mandrill"
	"github.com/gapidobri/prizer/internal/pkg/clients/sheets"
	"github.com/gapidobri/prizer/internal/pkg/models/config"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx := context.Background()

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.WithError(err).Fatal("Failed to decode config")
	}

	db, err := sqlx.Connect("postgres", cfg.Database.ConnectionString)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	// Clients
	addressValidationClient, err := addressvalidation.NewClient(ctx, cfg.AddressValidation.ApiKey)
	if err != nil {
		log.WithError(err).Fatal("Failed to create address validation client")
	}

	mandrillClient, err := mandrill.NewClient(cfg.Mandrill.ApiKey)
	if err != nil {
		log.WithError(err).Fatal("Failed to create mandrill client")
	}

	sheetsClient, err := sheets.NewClient(ctx, cfg.GoogleSheets.ServiceAccountKeyPath)
	if err != nil {
		log.WithError(err).Fatal("Failed to create sheets client")
	}

	// Repositories
	gameRepository := database.NewGameRepository(db)
	prizeRepository := database.NewPrizeRepository(db)
	wonPrizeRepository := database.NewWonPrizeRepository(db)
	drawMethodRepository := database.NewDrawMethodRepository(db)
	userRepository := database.NewUserRepository(db)
	participationMethodRepository := database.NewParticipationMethodRepository(db)
	participationRepository := database.NewParticipationRepository(db)
	mailTemplateRepository := database.NewMailTemplateRepository(db)

	// Services
	gameService := service.NewGameService(
		gameRepository,
		prizeRepository,
		wonPrizeRepository,
		userRepository,
		drawMethodRepository,
		participationMethodRepository,
		participationRepository,
		mailTemplateRepository,
		addressValidationClient,
		mandrillClient,
		sheetsClient,
	)
	userService := service.NewUserService(userRepository)
	prizeService := service.NewPrizeService(prizeRepository)
	wonPrizeService := service.NewWonPrizeService(wonPrizeRepository)
	participationMethodService := service.NewParticipationMethodService(participationMethodRepository)

	// APIs
	publicApi := public.NewServer(db, gameService)
	adminApi := admin.NewServer(
		db,
		gameService,
		userService,
		prizeService,
		wonPrizeService,
		participationMethodService,
	)

	go publicApi.Run(cfg.Http.Public.Address)
	go adminApi.Run(cfg.Http.Admin.Address)

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	log.Info("Shutting down...")
}
