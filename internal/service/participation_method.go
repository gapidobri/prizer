package service

import (
	"context"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/samber/lo"
)

type ParticipationMethodService struct {
	participationMethodRepository database.ParticipationMethodRepository
}

func NewParticipationMethodService(
	participationMethodRepository database.ParticipationMethodRepository,
) *ParticipationMethodService {
	return &ParticipationMethodService{
		participationMethodRepository: participationMethodRepository,
	}
}

func (s *ParticipationMethodService) GetParticipationMethods(
	ctx context.Context,
	filter api.GetParticipationMethodsFilter,
) (api.GetParticipationMethodsResponse, error) {
	participationMethods, err := s.participationMethodRepository.GetParticipationMethods(ctx, filter.ToDB())
	if err != nil {
		return nil, err
	}

	apiParticipationMethods := lo.Map(participationMethods, func(participationMethod dbModels.ParticipationMethod, _ int) api.ParticipationMethod {
		return api.ParticipationMethodFromDB(participationMethod)
	})
	return apiParticipationMethods, nil
}
