package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

type PublicPrize struct {
	// required: true
	Name string `json:"name"`

	Description *string `json:"description"`

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
	GameId string `json:"game_id"`

	// required: true
	Name string `json:"name"`

	Description *string `json:"description"`

	ImageUrl *string `json:"image_url"`

	// required: true
	Count int `json:"count"`

	// required: true
	WonCount int `json:"won_count"`
}

func PrizeFromDB(prize dbModels.Prize) Prize {
	return Prize{
		Id:          prize.Id,
		GameId:      prize.GameId,
		Name:        prize.Name,
		Description: prize.Description,
		ImageUrl:    prize.ImageUrl,
		Count:       prize.Count,
		WonCount:    prize.WonCount,
	}
}

// swagger:model GetPrizesResponse
type GetPrizesResponse []Prize

type GetPrizesFilter struct {
	GameId *string `form:"game_id" binding:"omitnil,uuid"`
}

func (f GetPrizesFilter) ToDB() dbModels.GetPrizesFilter {
	return dbModels.GetPrizesFilter{
		GameId: f.GameId,
	}
}

// swagger:model UpdatePrizeRequest
type UpdatePrizeRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"image_url" binding:"omitempty,url"`
	Count       int     `json:"count" binding:"required,number"`
}

func (r UpdatePrizeRequest) ToDB() dbModels.UpdatePrize {
	return dbModels.UpdatePrize{
		Name:        r.Name,
		Description: r.Description,
		ImageUrl:    r.ImageUrl,
		Count:       r.Count,
	}
}

// swagger:model CreatePrizeRequest
type CreatePrizeRequest struct {
	// required: true
	GameId string `json:"game_id" binding:"required,uuid"`

	// required: true
	Name string `json:"name" binding:"required"`

	Description *string `json:"description"`

	ImageUrl *string `json:"image_url" binding:"omitempty,url"`

	// required: true
	Count int `json:"count" binding:"required,number"`
}

func (r CreatePrizeRequest) ToDB() dbModels.CreatePrize {
	return dbModels.CreatePrize{
		GameId:      r.GameId,
		Name:        r.Name,
		Description: r.Description,
		ImageUrl:    r.ImageUrl,
		Count:       r.Count,
	}
}
