package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

type Prize struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func PrizeFromDB(prize dbModels.Prize) Prize {
	return Prize{
		Name:        prize.Name,
		Description: prize.Description,
	}
}
