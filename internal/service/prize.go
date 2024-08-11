package service

import (
	"context"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/samber/lo"
)

type PrizeService struct {
	prizeRepository database.PrizeRepository
}

func NewPrizeService(prizeRepository database.PrizeRepository) *PrizeService {
	return &PrizeService{
		prizeRepository: prizeRepository,
	}
}

func (s *PrizeService) GetPrizes(ctx context.Context, filter api.GetPrizesFilter) (api.GetPrizesResponse, error) {
	prizes, err := s.prizeRepository.GetPrizes(ctx, filter.ToDB())
	if err != nil {
		return nil, err
	}
	apiPrizes := lo.Map(prizes, func(prize dbModels.Prize, _ int) api.Prize {
		return api.PrizeFromDB(prize)
	})
	return apiPrizes, nil
}
