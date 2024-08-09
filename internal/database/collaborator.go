package database

import (
	"context"
	"database/sql"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type CollaboratorRepository interface {
	GetCollaboratorFromEmail(ctx context.Context, gameId string, email string) (*database.Collaborator, error)
	GetCollaboratorFromEmailAndAddress(ctx context.Context, gameId string, email string, address string) (*database.Collaborator, error)
	CreateCollaborator(ctx context.Context, collaborator database.CreateCollaborator, unique bool) (*database.Collaborator, error)
	UpdateLastRollTime(ctx context.Context, collaboratorId string, time time.Time) error
}

type collaboratorRepository struct {
	db *sqlx.DB
}

func NewCollaboratorRepository(db *sqlx.DB) CollaboratorRepository {
	return &collaboratorRepository{
		db: db,
	}
}

func (c *collaboratorRepository) GetCollaboratorFromEmail(ctx context.Context, gameId string, email string) (*database.Collaborator, error) {
	var collaborator database.Collaborator
	err := c.db.GetContext(ctx, &collaborator, `
		SELECT *
		FROM collaborator
		WHERE email = $1
			AND game_id = $2
	`, email, gameId)
	if err != nil {
		return nil, err
	}
	return &collaborator, nil
}

func (c *collaboratorRepository) GetCollaboratorFromEmailAndAddress(ctx context.Context, gameId string, email string, address string) (*database.Collaborator, error) {
	var collaborator database.Collaborator
	err := c.db.GetContext(ctx, &collaborator, `
		SELECT *
		FROM collaborator
		WHERE email = $1
	    	AND address = $2
			AND game_id = $3
	`, email, address, gameId)
	if err != nil {
		return nil, err
	}
	return &collaborator, nil
}

func (c *collaboratorRepository) CreateCollaborator(ctx context.Context, create database.CreateCollaborator, unique bool) (*database.Collaborator, error) {
	if unique {
		var collaborator database.Collaborator
		err := c.db.GetContext(ctx, &collaborator, `
			SELECT *
			FROM collaborator
			WHERE game_id = $1
				AND (email = $2 OR address = $3)
		`)
		switch {
		case err == nil:
			return nil, er.CollaboratorExists
		case errors.Is(err, sql.ErrNoRows):
			break
		default:
			return nil, err
		}
	}

	_, err := c.db.NamedExecContext(ctx, `
		INSERT INTO collaborator (email, game_id)
		VALUES (:email, :game_id)
	`, create)
	if err != nil {
		return nil, err
	}

	collaborator, err := c.GetCollaboratorFromEmail(ctx, create.Email, create.GameId)
	if err != nil {
		return nil, err
	}

	return collaborator, nil
}

func (c *collaboratorRepository) UpdateLastRollTime(ctx context.Context, collaboratorId string, time time.Time) error {
	_, err := c.db.ExecContext(ctx, `
		UPDATE collaborator
		SET last_roll_time = $1
		WHERE collaborator_id = $2
	`, time, collaboratorId)
	if err != nil {
		return err
	}
	return nil
}
