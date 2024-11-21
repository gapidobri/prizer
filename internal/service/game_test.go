package service

import (
	"context"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/clients/addressvalidation"
	"github.com/gapidobri/prizer/internal/pkg/clients/mandrill"
	"github.com/gapidobri/prizer/internal/pkg/clients/sheets"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	game = dbModels.Game{
		Id: "049fce3d-46d6-4dbf-ad70-3b93e98e6d58",
	}
	participationMethod = dbModels.ParticipationMethod{
		Id:                 "3ef7669-0179-4347-895a-cb8b045f74f8",
		GameId:             game.Id,
		ParticipationLimit: enums.ParticipationLimitNone,
		Fields: dbModels.FieldConfig{
			User: map[string]dbModels.Field{
				"name": {
					Type:     enums.FieldTypeString,
					Required: true,
					Unique:   false,
				},
				"email": {
					Type:     enums.FieldTypeString,
					Required: true,
					Unique:   true,
				},
			},
			Participation: map[string]dbModels.Field{
				"code": {
					Type:     enums.FieldTypeString,
					Required: true,
					Unique:   true,
				},
			},
		},
	}
	drawMethodChanceAlways = dbModels.DrawMethod{
		Id:     "4af50df6-67a4-4771-80f3-a5267f9af684",
		Method: dbModels.DrawMethodChance,
		Data:   "{\"chance\": \"1\"}",
	}
	drawMethodChanceNever = dbModels.DrawMethod{
		Id:     "0a4079bd-6fbe-4e08-ae97-7ba993b81a30",
		Method: dbModels.DrawMethodChance,
		Data:   "{\"chance\": \"0\"}",
	}
	drawMethodFirstN = dbModels.DrawMethod{
		Id:     "e8840492-7ece-4d90-a0b2-02ebea5752b3",
		Method: dbModels.DrawMethodFirstN,
		Data:   "{}",
	}
	user = dbModels.User{
		Id:     "564a1991-caff-44fd-98f2-12e64c30aa8c",
		GameId: game.Id,
		UserFields: dbModels.UserFields{
			Email: lo.ToPtr("test@example.com"),
		},
		AdditionalFields: map[string]any{
			"name": "Test User",
		},
	}
	prize = dbModels.Prize{
		Id:          "560b6b8b-2ea0-4bf8-b38c-425ff3175217",
		GameId:      game.Id,
		Name:        "Prize Name",
		Description: "Prize Description",
		Count:       1,
	}
	participation = dbModels.Participation{
		Id:                    "3638e1e5-894d-4df5-947e-61e04e046d6e",
		UserId:                user.Id,
		ParticipationMethodId: participationMethod.Id,
	}
)

func TestGameService_Participate(t *testing.T) {
	ctx := context.Background()

	gameRepo := database.NewGameRepositoryMock(t)
	prizeRepo := database.NewPrizeRepositoryMock(t)
	wonPrizeRepo := database.NewWonPrizeRepositoryMock(t)
	userRepo := database.NewUserRepositoryMock(t)
	drawMethodRepo := database.NewDrawMethodRepositoryMock(t)
	participationMethodRepo := database.NewParticipationMethodRepositoryMock(t)
	participationRepo := database.NewParticipationRepositoryMock(t)
	mailTemplateRepo := database.NewMailTemplateRepositoryMock(t)
	addressValidationClient := addressvalidation.NewClientMock(t)
	mandrillClient := mandrill.NewClientMock(t)
	sheetsClient := sheets.NewClientMock(t)

	gameService := NewGameService(gameRepo, prizeRepo, wonPrizeRepo, userRepo, drawMethodRepo, participationMethodRepo,
		participationRepo, mailTemplateRepo, addressValidationClient, mandrillClient, sheetsClient)

	// Function mocks

	participationMethodRepo.OnGetParticipationMethod(participationMethod.Id).TypedReturns(&participationMethod, nil)

	gameRepo.OnGetGame(game.Id).TypedReturns(&game, nil)

	userRepo.OnGetUserFromFields(game.Id, user.UserFields).TypedReturns(&user, nil)

	participationRepo.OnGetParticipations(dbModels.GetParticipationsFilter{
		UserId:                &user.Id,
		ParticipationMethodId: &participationMethod.Id,
		Fields: &dbModels.JsonMap{
			"code": "12345",
		},
	}).TypedReturns([]dbModels.Participation{}, nil)

	drawMethodRepo.OnGetDrawMethods(game.Id, dbModels.GetDrawMethodsFilter{
		ParticipationMethodId: &participationMethod.Id,
	}).TypedReturns([]dbModels.DrawMethod{drawMethodFirstN}, nil)

	prizeRepo.OnGetPrizes(dbModels.GetPrizesFilter{
		GameId:        &game.Id,
		DrawMethodId:  &drawMethodFirstN.Id,
		AvailableOnly: true,
	}).TypedReturns([]dbModels.Prize{prize}, nil)

	participationRepo.OnCreateParticipation(dbModels.CreateParticipation{
		UserId:                user.Id,
		ParticipationMethodId: participationMethod.Id,
		Fields: dbModels.JsonMap{
			"code": "12345",
		},
	}).TypedReturns(&participation, nil)

	wonPrizeRepo.OnCreateWonPrize(dbModels.CreateWonPrize{
		PrizeId:         prize.Id,
		ParticipationId: participation.Id,
	}).TypedReturns(nil)

	// Test

	request := api.ParticipationRequest{
		Fields: map[string]any{
			"email": "test@example.com",
			"name":  "Test User",
			"code":  "12345",
		},
	}

	response, err := gameService.Participate(ctx, participationMethod.Id, request)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, response, &api.ParticipationResponse{
		Prizes: []api.PublicPrize{api.PublicPrizeFromDB(prize)},
	})
}
