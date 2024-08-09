package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/client/addressvalidation"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/util"
	"math/rand"
	"time"
)

type GameService struct {
	gameRepository          database.GameRepository
	prizeRepository         database.PrizeRepository
	wonPrizeRepository      database.WonPrizeRepository
	collaboratorRepository  database.CollaboratorRepository
	addressValidationClient *addressvalidation.Client
}

func NewGameService(
	gameRepository database.GameRepository,
	prizeRepository database.PrizeRepository,
	wonPrizeRepository database.WonPrizeRepository,
	collaboratorRepository database.CollaboratorRepository,
	addressValidationClient *addressvalidation.Client,
) *GameService {
	return &GameService{
		gameRepository:          gameRepository,
		prizeRepository:         prizeRepository,
		wonPrizeRepository:      wonPrizeRepository,
		collaboratorRepository:  collaboratorRepository,
		addressValidationClient: addressValidationClient,
	}
}

func (g *GameService) Roll(ctx context.Context, gameId string, roll api.RollRequest) (*api.RollResponse, error) {
	game, err := g.gameRepository.GetGame(ctx, gameId)
	if err != nil {
		return nil, err
	}

	normalizedAddress, err := g.addressValidationClient.NormalizeAddress(ctx, roll.Address)
	if err != nil {
		return nil, err
	}

	collaborator, err := g.collaboratorRepository.GetCollaboratorFromEmailAndAddress(ctx, gameId, roll.Email, normalizedAddress)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		collaborator, err = g.collaboratorRepository.CreateCollaborator(ctx, dbModels.CreateCollaborator{
			GameId:  game.Id,
			Email:   roll.Email,
			Address: &normalizedAddress,
		}, game.UniqueCollaboratorData)
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

	if game.WinPercentage <= rand.Float32() {
		return loose()
	}

	prizes, err := g.prizeRepository.GetPrizes(ctx, dbModels.GetPrizesFilter{GameId: &gameId, AvailableOnly: true})
	if err != nil {
		return nil, err
	}

	if len(prizes) == 0 {
		return loose()
	}

	// TODO: implement different chances

	prize := prizes[rand.Intn(len(prizes))]

	err = g.wonPrizeRepository.CreateWonPrize(ctx, dbModels.CreateWonPrize{
		PrizeId:        prize.Id,
		CollaboratorId: collaborator.Id,
	})
	if err != nil {
		return nil, err
	}

	return &api.RollResponse{
		Won: true,
		Prize: &api.Prize{
			Name:        prize.Name,
			Description: prize.Description,
		},
	}, nil
}
