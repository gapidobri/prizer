package internal

import (
	"context"
	"github.com/gapidobri/prizer/internal/api/admin"
	"github.com/gapidobri/prizer/internal/api/public"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/client/addressvalidation"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx := context.Background()

	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Clients
	addressValidationClient, err := addressvalidation.NewClient(ctx, viper.GetString("address_validation_api_key"))
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	gameRepository := database.NewGameRepository(db)
	prizeRepository := database.NewPrizeRepository(db)
	wonPrizeRepository := database.NewWonPrizeRepository(db)
	drawMethodRepository := database.NewDrawMethodRepository(db)
	userRepository := database.NewUserRepository(db)
	participationMethodRepository := database.NewParticipationMethodRepository(db)
	participationRepository := database.NewParticipationRepository(db)

	// Services
	gameService := service.NewGameService(
		gameRepository,
		prizeRepository,
		wonPrizeRepository,
		userRepository,
		drawMethodRepository,
		participationMethodRepository,
		participationRepository,
		addressValidationClient,
	)
	userService := service.NewUserService(userRepository)
	prizeService := service.NewPrizeService(prizeRepository)
	wonPrizeService := service.NewWonPrizeService(wonPrizeRepository)
	participationMethodService := service.NewParticipationMethodService(participationMethodRepository)

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
