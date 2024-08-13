package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

// swagger:model Game
type Game struct {
	// required: true
	Id string `json:"id"`

	// required: true
	Name string `json:"name"`

	GoogleSheetId *string `json:"google_sheet_id"`

	GoogleSheetTabName *string `json:"google_sheet_tab_name"`
}

func GameFromDB(game dbModels.Game) Game {
	return Game{
		Id:                 game.Id,
		Name:               game.Name,
		GoogleSheetId:      game.GoogleSheetId,
		GoogleSheetTabName: game.GoogleSheetTabName,
	}
}

// swagger:model GetGamesResponse
type GetGamesResponse []Game

// swagger:model GetGameResponse
type GetGameResponse *Game
