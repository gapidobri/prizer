package service

import (
	"context"
	"github.com/gapidobri/prizer/internal/database"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/google/uuid"
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

func (s *PrizeService) CreatePrize(ctx context.Context, prize api.CreatePrizeRequest) error {
	return s.prizeRepository.CreatePrize(ctx, prize.ToDB())
}

func (s *PrizeService) UpdatePrize(ctx context.Context, prizeId string, prize api.UpdatePrizeRequest) error {
	err := uuid.Validate(prizeId)
	if err != nil {
		return er.InvalidUuid
	}
	return s.prizeRepository.UpdatePrize(ctx, prizeId, prize.ToDB())
}

func (s *PrizeService) DeletePrize(ctx context.Context, prizeId string) error {
	err := uuid.Validate(prizeId)
	if err != nil {
		return er.InvalidUuid
	}
	return s.prizeRepository.DeletePrize(ctx, prizeId)
}
