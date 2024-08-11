package service

import (
	"context"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/samber/lo"
)

type WonPrizeService struct {
	wonPrizeRepository database.WonPrizeRepository
}

func NewWonPrizeService(wonPrizeRepository database.WonPrizeRepository) *WonPrizeService {
	return &WonPrizeService{
		wonPrizeRepository: wonPrizeRepository,
	}
}

func (s *WonPrizeService) GetWonPrizes(ctx context.Context, filter api.GetWonPrizesFilter) (api.GetWonPrizesResponse, error) {
	wonPrizes, err := s.wonPrizeRepository.GetWonPrizes(ctx, filter.ToDB())
	if err != nil {
		return nil, err
	}

	apiWonPrizes := lo.Map(wonPrizes, func(wonPrize dbModels.WonPrize, _ int) api.WonPrize {
		return api.WonPrizeFromDB(wonPrize)
	})
	return apiWonPrizes, nil
}
