package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

// swagger:model WonPrize
type WonPrize struct {
	// required: true
	Prize Prize `json:"prize"`

	// required: true
	Participation Participation `json:"participation"`

	// required: true
	User User `json:"user"`
}

func WonPrizeFromDB(wonPrize dbModels.WonPrize) WonPrize {
	return WonPrize{
		Prize:         PrizeFromDB(wonPrize.Prize),
		User:          UserFromDB(wonPrize.User),
		Participation: ParticipationFromDB(wonPrize.Participation),
	}
}

type GetWonPrizesFilter struct {
	GameId  *string `form:"gameId"`
	UserId  *string `form:"userId"`
	PrizeId *string `form:"prizeId"`
}

func (f GetWonPrizesFilter) ToDB() dbModels.GetWonPrizesFilter {
	return dbModels.GetWonPrizesFilter{
		GameId:  f.GameId,
		UserId:  f.UserId,
		PrizeId: f.PrizeId,
	}
}

// swagger:model GetWonPrizesResponse
type GetWonPrizesResponse []WonPrize
