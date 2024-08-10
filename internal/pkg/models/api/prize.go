package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

type PublicPrize struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func PublicPrizeFromDB(prize dbModels.Prize) PublicPrize {
	return PublicPrize{
		Name:        prize.Name,
		Description: prize.Description,
	}
}

type Prize struct {
	Id          string `json:"id"`
	GameId      string `json:"gameId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

func PrizeFromDB(prize dbModels.Prize) Prize {
	return Prize{
		Id:          prize.Id,
		GameId:      prize.GameId,
		Name:        prize.Name,
		Description: prize.Description,
		Count:       prize.Count,
	}
}

// swagger:model GetPrizesResponse
type GetPrizesResponse []Prize
