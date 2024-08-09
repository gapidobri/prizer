package database

import (
	"context"
	"database/sql"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type GameRepository interface {
	GetGames(ctx context.Context) ([]database.Game, error)
	GetGame(ctx context.Context, gameId string) (*database.Game, error)
}

type gameRepository struct {
	db *sqlx.DB
}

func NewGameRepository(db *sqlx.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

func (g *gameRepository) GetGames(ctx context.Context) ([]database.Game, error) {
	var games []database.Game
	err := g.db.SelectContext(ctx, &games, "SELECT * FROM game")
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (g *gameRepository) GetGame(ctx context.Context, gameId string) (*database.Game, error) {
	var game database.Game
	err := g.db.GetContext(ctx, &game, `
		SELECT *
		FROM game g
		WHERE g.game_id = $1
	`, gameId)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		return nil, er.GameNotFound
	default:
		return nil, err
	}

	return &game, nil
}
