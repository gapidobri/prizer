package internal

import (
	"context"
	"github.com/gapidobri/prizer/internal/api/admin"
	"github.com/gapidobri/prizer/internal/api/public"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/client/addressvalidation"
	"github.com/gapidobri/prizer/internal/pkg/models/config"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattbaird/gochimp"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx := context.Background()

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to decode config, %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("failed to connect to database, %v", err)
	}

	// Clients
	addressValidationClient, err := addressvalidation.NewClient(ctx, cfg.AddressValidation.ApiKey)
	if err != nil {
		log.Fatal(err)
	}

	mandrillClient, err := gochimp.NewMandrill(cfg.Mandrill.ApiKey)
	if err != nil {
		log.Fatalf("failed to create mandrill client, %v", err)
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
	)
	userService := service.NewUserService(userRepository)
	prizeService := service.NewPrizeService(prizeRepository)
	wonPrizeService := service.NewWonPrizeService(wonPrizeRepository)
	participationMethodService := service.NewParticipationMethodService(participationMethodRepository)

	// APIs
	publicApi := public.NewServer(gameService)
	adminApi := admin.NewServer(
		gameService,
		userService,
		prizeService,
		wonPrizeService,
		participationMethodService,
	)

	go publicApi.Run(":8080")
	go adminApi.Run(":8081")

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	log.Println("Shutting down...")
}
