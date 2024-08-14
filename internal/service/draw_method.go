package service

import (
	"context"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/samber/lo"
)

type DrawMethodService struct {
	drawMethodRepository database.DrawMethodRepository
}

func NewDrawMethodService(drawMethodRepository database.DrawMethodRepository) *DrawMethodService {
	return &DrawMethodService{
		drawMethodRepository: drawMethodRepository,
	}
}

func (s *DrawMethodService) GetDrawMethods(ctx context.Context, filter api.GetDrawMethodsFilter) (api.GetDrawMethodsResponse, error) {
	drawMethods, err := s.drawMethodRepository.GetDrawMethods(ctx, filter.ToDB())
	if err != nil {
		return nil, err
	}

	apiDrawMethods := lo.Map(drawMethods, func(drawMethod dbModels.DrawMethod, _ int) api.DrawMethod {
		return api.DrawMethodFromDB(drawMethod)
	})

	return apiDrawMethods, nil
}
