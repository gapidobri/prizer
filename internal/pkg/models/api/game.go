package api

import (
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
)

// swagger:model RollRequest
type RollRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Address string `json:"address" binding:"required"`
}

// swagger:model RollResponse
type RollResponse struct {
	Won    bool    `json:"won"`
	Prizes []Prize `json:"prizes"`
}

// swagger:model Game
type Game struct {
	Id                     string  `json:"id"`
	Name                   string  `json:"name"`
	WinPercentage          float32 `json:"win_percentage"`
	UniqueCollaboratorData bool    `json:"unique_collaborator_data"`
}

func GameFromDB(game dbModels.Game) Game {
	return Game{
		Id:                     game.Id,
		Name:                   game.Name,
		WinPercentage:          game.WinPercentage,
		UniqueCollaboratorData: game.UniqueCollaboratorData,
	}
}

// swagger:model GetGamesResponse
type GetGamesResponse = []Game

// swagger:model GetGameResponse
type GetGameResponse = Game
