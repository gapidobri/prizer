package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/client/addressvalidation"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/util"
	"github.com/samber/lo"
	"math/rand"
	"time"
)

type GameService struct {
	gameRepository                database.GameRepository
	prizeRepository               database.PrizeRepository
	wonPrizeRepository            database.WonPrizeRepository
	collaboratorRepository        database.CollaboratorRepository
	drawMethodRepository          database.DrawMethodRepository
	collaborationMethodRepository database.CollaborationMethodRepository
	collaborationRepository       database.CollaborationRepository
	addressValidationClient       *addressvalidation.Client
}

func NewGameService(
	gameRepository database.GameRepository,
	prizeRepository database.PrizeRepository,
	wonPrizeRepository database.WonPrizeRepository,
	collaboratorRepository database.CollaboratorRepository,
	drawMethodRepository database.DrawMethodRepository,
	collaborationMethodRepository database.CollaborationMethodRepository,
	collaborationRepository database.CollaborationRepository,
	addressValidationClient *addressvalidation.Client,
) *GameService {
	return &GameService{
		gameRepository:                gameRepository,
		prizeRepository:               prizeRepository,
		wonPrizeRepository:            wonPrizeRepository,
		collaboratorRepository:        collaboratorRepository,
		drawMethodRepository:          drawMethodRepository,
		collaborationMethodRepository: collaborationMethodRepository,
		collaborationRepository:       collaborationRepository,
		addressValidationClient:       addressValidationClient,
	}
}

func (g *GameService) GetGames(ctx context.Context) ([]api.Game, error) {
	games, err := g.gameRepository.GetGames(ctx)
	if err != nil {
		return nil, err
	}

	apiGames := lo.Map(games, func(game dbModels.Game, index int) api.Game {
		return api.GameFromDB(game)
	})

	return apiGames, nil
}

func (g *GameService) GetGame(ctx context.Context, gameId string) (*api.Game, error) {
	game, err := g.gameRepository.GetGame(ctx, gameId)
	if err != nil {
		return nil, err
	}

	apiGame := api.GameFromDB(*game)

	return &apiGame, nil
}

func (g *GameService) Roll(ctx context.Context, collaborationMethodId string, roll api.RollRequest) (*api.RollResponse, error) {
	collaborationMethod, err := g.collaborationMethodRepository.GetCollaborationMethod(ctx, collaborationMethodId)
	if err != nil {
		return nil, err
	}

	normalizedAddress, err := g.addressValidationClient.NormalizeAddress(ctx, roll.Address)
	if err != nil {
		return nil, err
	}

	collaborator, err := g.collaboratorRepository.GetCollaboratorFromEmailAndAddress(ctx, collaborationMethod.GameId, roll.Email, normalizedAddress)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		collaborator, err = g.collaboratorRepository.CreateCollaborator(ctx, dbModels.CreateCollaborator{
			GameId:  collaborationMethod.GameId,
			Email:   roll.Email,
			Address: &normalizedAddress,
		}, false)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	if collaborator.LastRollTime != nil &&
		!util.StripTime(*collaborator.LastRollTime).Before(util.StripTime(time.Now())) {
		return nil, er.AlreadyRolled
	}

	loose := func() (*api.RollResponse, error) {
		err := g.collaboratorRepository.UpdateLastRollTime(ctx, collaborator.Id, time.Now())
		if err != nil {
			return nil, err
		}
		return &api.RollResponse{Won: false}, nil
	}

	err = g.collaboratorRepository.UpdateLastRollTime(ctx, collaborator.Id, time.Now())
	if err != nil {
		return nil, err
	}

	drawMethods, err := g.drawMethodRepository.GetDrawMethods(ctx, collaborationMethod.GameId, dbModels.GetDrawMethodsFilter{
		CollaborationMethodId: &collaborationMethodId,
	})
	if err != nil {
		return nil, err
	}

	var wonPrizes []dbModels.Prize

	for _, drawMethod := range drawMethods {
		var prizes []dbModels.Prize
		prizes, err = g.prizeRepository.GetPrizes(ctx, dbModels.GetPrizesFilter{
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

	if len(wonPrizes) == 0 {
		return loose()
	}

	for _, prize := range wonPrizes {
		err = g.wonPrizeRepository.CreateWonPrize(ctx, dbModels.CreateWonPrize{
			PrizeId:        prize.Id,
			CollaboratorId: collaborator.Id,
		})
		if err != nil {
			return nil, err
		}
	}

	apiPrizes := lo.Map(wonPrizes, func(prize dbModels.Prize, _ int) api.Prize {
		return api.PrizeFromDB(prize)
	})

	return &api.RollResponse{
		Won:    true,
		Prizes: apiPrizes,
	}, nil
}
