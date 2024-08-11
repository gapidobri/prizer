package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/client/addressvalidation"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
	"github.com/gapidobri/prizer/internal/pkg/util"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"math/rand"
	"net/mail"
	"time"
)

type GameService struct {
	gameRepository                database.GameRepository
	prizeRepository               database.PrizeRepository
	wonPrizeRepository            database.WonPrizeRepository
	userRepository                database.UserRepository
	drawMethodRepository          database.DrawMethodRepository
	participationMethodRepository database.ParticipationMethodRepository
	participationRepository       database.ParticipationRepository
	addressValidationClient       *addressvalidation.Client
}

func NewGameService(
	gameRepository database.GameRepository,
	prizeRepository database.PrizeRepository,
	wonPrizeRepository database.WonPrizeRepository,
	userRepository database.UserRepository,
	drawMethodRepository database.DrawMethodRepository,
	participationMethodRepository database.ParticipationMethodRepository,
	participationRepository database.ParticipationRepository,
	addressValidationClient *addressvalidation.Client,
) *GameService {
	return &GameService{
		gameRepository:                gameRepository,
		prizeRepository:               prizeRepository,
		wonPrizeRepository:            wonPrizeRepository,
		userRepository:                userRepository,
		drawMethodRepository:          drawMethodRepository,
		participationMethodRepository: participationMethodRepository,
		participationRepository:       participationRepository,
		addressValidationClient:       addressValidationClient,
	}
}

func (s *GameService) GetGames(ctx context.Context) (api.GetGamesResponse, error) {
	games, err := s.gameRepository.GetGames(ctx)
	if err != nil {
		return nil, err
	}

	apiGames := lo.Map(games, func(game dbModels.Game, index int) api.Game {
		return api.GameFromDB(game)
	})

	return apiGames, nil
}

func (s *GameService) GetGame(ctx context.Context, gameId string) (api.GetGameResponse, error) {
	err := uuid.Validate(gameId)
	if err != nil {
		return nil, er.InvalidUuid
	}

	game, err := s.gameRepository.GetGame(ctx, gameId)
	if err != nil {
		return nil, err
	}

	apiGame := api.GameFromDB(*game)

	return &apiGame, nil
}

func (s *GameService) Participate(ctx context.Context, participationMethodId string, roll api.ParticipationRequest) (*api.ParticipationResponse, error) {
	participationMethod, err := s.participationMethodRepository.GetParticipationMethod(ctx, participationMethodId)
	if err != nil {
		return nil, err
	}

	gameId := participationMethod.GameId

	// Sanitize and validate incoming data
	additionalUserFields := dbModels.JsonMap{}
	var (
		userFields   dbModels.UserFields
		uniqueFields dbModels.UserFields
	)
	for key, field := range participationMethod.Fields.User {
		value, exists := roll.Fields[key]
		if !exists {
			if field.Unique || field.Required {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is required", key))
			}
			continue
		}

		switch key {
		case "email":
			str, ok := value.(string)
			if !ok {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
			}
			email, err := mail.ParseAddress(str)
			if err != nil {
				return nil, er.InvalidEmail
			}

			userFields.Email = &email.Address
			if field.Unique {
				uniqueFields.Email = &email.Address
			}
			continue

		case "address":
			str, ok := value.(string)
			if !ok {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
			}

			address, err := s.addressValidationClient.NormalizeAddress(ctx, str)
			if err != nil {
				return nil, err
			}

			userFields.Address = &address
			if field.Unique {
				uniqueFields.Address = &address
			}
			continue
		}

		value, err = s.validateField(field, key, value)
		if err != nil {
			return nil, err
		}
		additionalUserFields[key] = value
	}

	participationFields := dbModels.JsonMap{}
	uniqueParticipationFields := dbModels.JsonMap{}
	for key, field := range participationMethod.Fields.Participation {
		value, exists := roll.Fields[key]
		if !exists {
			if field.Unique || field.Required {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is required", key))
			}
			continue
		}

		value, err = s.validateField(field, key, value)
		if err != nil {
			return nil, err
		}
		participationFields[key] = value
		if field.Unique {
			uniqueParticipationFields[key] = value
		}
	}

	// Get / create user
	user, err := s.userRepository.GetUserFromFields(ctx, gameId, uniqueFields)
	switch {
	case err == nil:
		break
	case errors.Is(err, er.UserNotFound):
		user, err = s.userRepository.CreateUser(ctx, dbModels.CreateUser{
			GameId:           gameId,
			UserFields:       userFields,
			AdditionalFields: additionalUserFields,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	// Check if user can participate
	switch participationMethod.Limit {
	case enums.ParticipationLimitNone:
		break
	case enums.ParticipationLimitDaily:
		participations, err := s.participationRepository.GetParticipations(ctx, dbModels.GetParticipationsFilter{
			UserId:                &user.Id,
			ParticipationMethodId: &participationMethod.Id,
			From:                  lo.ToPtr(util.StripTime(time.Now())),
			To:                    lo.ToPtr(time.Now()),
		})
		if err != nil {
			return nil, err
		}
		if len(participations) != 0 {
			return nil, er.AlreadyParticipated
		}
	}

	// Check participation unique fields
	if len(uniqueParticipationFields) > 0 {
		participations, err := s.participationRepository.GetParticipations(ctx, dbModels.GetParticipationsFilter{
			UserId:                &user.Id,
			ParticipationMethodId: &participationMethod.Id,
			Fields:                &uniqueParticipationFields,
		})
		if err != nil {
			return nil, err
		}
		if len(participations) != 0 {
			return nil, er.ParticipationDataExists
		}
	}

	// Participate in all draw methods
	drawMethods, err := s.drawMethodRepository.GetDrawMethods(ctx, gameId, dbModels.GetDrawMethodsFilter{
		ParticipationMethodId: &participationMethodId,
	})
	if err != nil {
		return nil, err
	}

	var wonPrizes []dbModels.Prize

	for _, drawMethod := range drawMethods {
		var prizes []dbModels.Prize
		prizes, err = s.prizeRepository.GetPrizes(ctx, dbModels.GetPrizesFilter{
			GameId:        &gameId,
			DrawMethodId:  &drawMethod.Id,
			AvailableOnly: true,
		})
		if err != nil {
			return nil, err
		}

		if len(prizes) == 0 {
			continue
		}

		switch drawMethod.Method {
		case dbModels.DrawMethodFirstN:
			wonPrizes = append(wonPrizes, prizes[0])

		case dbModels.DrawMethodChance:
			var data dbModels.DrawMethodChanceData
			err = json.Unmarshal([]byte(drawMethod.Data), &data)
			if err != nil {
				return nil, err
			}

			if data.Chance <= rand.Float64() {
				continue
			}

			prize := prizes[rand.Intn(len(prizes))]

			wonPrizes = append(wonPrizes, prize)
		}
	}

	participation, err := s.participationRepository.CreateParticipation(ctx, dbModels.CreateParticipation{
		ParticipationMethodId: participationMethodId,
		UserId:                user.Id,
		Fields:                participationFields,
	})
	if err != nil {
		return nil, err
	}

	for _, prize := range wonPrizes {
		err = s.wonPrizeRepository.CreateWonPrize(ctx, dbModels.CreateWonPrize{
			ParticipationId: participation.Id,
			PrizeId:         prize.Id,
		})
		if err != nil {
			return nil, err
		}
	}

	publicPrizes := lo.Map(wonPrizes, func(prize dbModels.Prize, _ int) api.PublicPrize {
		return api.PublicPrizeFromDB(prize)
	})

	return &api.ParticipationResponse{
		Prizes: publicPrizes,
	}, nil
}

func (s *GameService) validateField(field dbModels.Field, key string, value any) (any, error) {
	switch field.Type {
	case enums.FieldTypeBool:
		val, ok := value.(bool)
		if !ok {
			return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a bool", key))
		}
		return val, nil

	case enums.FieldTypeString:
		val, ok := value.(string)
		if !ok {
			return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
		}
		return val, nil

	default:
		return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a valid field", key))
	}
}
