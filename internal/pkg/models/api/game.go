package api

import (
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
)

// swagger:model Game
type Game struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GameFromDB(game dbModels.Game) Game {
	return Game{
		Id:   game.Id,
		Name: game.Name,
	}
}

// swagger:model GetGamesResponse
type GetGamesResponse = []Game

// swagger:model GetGameResponse
type GetGameResponse = *Game
