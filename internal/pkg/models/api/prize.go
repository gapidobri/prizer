package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

type PublicPrize struct {
	// required: true
	Name string `json:"name"`

	// required: true
	Description string `json:"description"`

	ImageUrl *string `json:"image_url"`
}

func PublicPrizeFromDB(prize dbModels.Prize) PublicPrize {
	return PublicPrize{
		Name:        prize.Name,
		Description: prize.Description,
		ImageUrl:    prize.ImageUrl,
	}
}

// swagger:model Prize
type Prize struct {
	// required: true
	Id string `json:"id"`

	// required: true
	GameId string `json:"gameId"`

	// required: true
	Name string `json:"name"`

	// required: true
	Description string `json:"description"`

	ImageUrl *string `json:"image_url"`

	// required: true
	Count int `json:"count"`
}

func PrizeFromDB(prize dbModels.Prize) Prize {
	return Prize{
		Id:          prize.Id,
		GameId:      prize.GameId,
		Name:        prize.Name,
		Description: prize.Description,
		ImageUrl:    prize.ImageUrl,
		Count:       prize.Count,
	}
}

// swagger:model GetPrizesResponse
type GetPrizesResponse []Prize

type GetPrizesFilter struct {
	GameId *string `form:"gameId" binding:"omitnil,uuid"`
}

func (f GetPrizesFilter) ToDB() dbModels.GetPrizesFilter {
	return dbModels.GetPrizesFilter{
		GameId: f.GameId,
	}
}
