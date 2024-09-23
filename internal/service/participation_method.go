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

func (s *ParticipationMethodService) UpdateParticipationMethod(
	ctx context.Context,
	participationMethodId string,
	participationMethod api.UpdateParticipationMethodRequest,
) error {
	err := uuid.Validate(participationMethodId)
	if err != nil {
		return er.InvalidUuid
	}

	return s.participationMethodRepository.UpdateParticipationMethod(ctx, participationMethodId, participationMethod.ToDB())
}

func (s *ParticipationMethodService) LinkDrawMethod(ctx context.Context, participationMethodId string, drawMethodId string) error {
	err := uuid.Validate(participationMethodId)
	if err != nil {
		return er.InvalidUuid
	}
	err = uuid.Validate(drawMethodId)
	if err != nil {
		return er.InvalidUuid
	}
	return s.LinkDrawMethod(ctx, participationMethodId, drawMethodId)
}

func (s *ParticipationMethodService) UnlinkDrawMethod(ctx context.Context, participationMethodId string, drawMethodId string) error {
	err := uuid.Validate(participationMethodId)
	if err != nil {
		return er.InvalidUuid
	}
	err = uuid.Validate(drawMethodId)
	if err != nil {
		return er.InvalidUuid
	}
	return s.UnlinkDrawMethod(ctx, participationMethodId, drawMethodId)
}
