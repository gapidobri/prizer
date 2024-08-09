package internal

import (
	"context"
	"github.com/gapidobri/prizer/internal/api"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/client/addressvalidation"
	"github.com/gapidobri/prizer/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func Run() {
	ctx := context.Background()

	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Clients
	addressValidationClient, err := addressvalidation.NewClient(ctx, "<access_token>")
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	gameRepository := database.NewGameRepository(db)
	prizeRepository := database.NewPrizeRepository(db)
	wonPrizeRepository := database.NewWonPrizeRepository(db)
	collaboratorRepository := database.NewCollaboratorRepository(db)

	// Services
	gameService := service.NewGameService(
		gameRepository,
		prizeRepository,
		wonPrizeRepository,
		collaboratorRepository,
		addressValidationClient,
	)

	server := api.NewServer(gameService)

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
