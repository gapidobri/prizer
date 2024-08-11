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

type UserService struct {
	userRepository database.UserRepository
}

func NewUserService(userRepository database.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUsers(ctx context.Context, filter api.GetUsersFilter) (api.GetUsersResponse, error) {
	users, err := s.userRepository.GetUsers(ctx, filter.ToDB())
	if err != nil {
		return nil, err
	}

	apiUsers := lo.Map(users, func(user dbModels.User, _ int) api.User {
		return api.UserFromDB(user)
	})

	return apiUsers, nil
}

func (s *UserService) GetUser(ctx context.Context, userId string) (api.GetUserResponse, error) {
	err := uuid.Validate(userId)
	if err != nil {
		return nil, er.InvalidUuid
	}

	user, err := s.userRepository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	apiUser := api.UserFromDB(*user)
	return &apiUser, nil
}
