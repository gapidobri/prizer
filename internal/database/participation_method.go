package database

import (
	"context"
	"database/sql"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ParticipationMethodRepository interface {
	GetParticipationMethod(ctx context.Context, participationMethodId string) (*database.ParticipationMethod, error)
}

type participationMethodRepository struct {
	db *sqlx.DB
}

func NewParticipationMethodRepository(db *sqlx.DB) ParticipationMethodRepository {
	return &participationMethodRepository{
		db: db,
	}
}

func (c *participationMethodRepository) GetParticipationMethod(ctx context.Context, participationMethodId string) (*database.ParticipationMethod, error) {
	var participationMethod database.ParticipationMethod
	err := c.db.GetContext(ctx, &participationMethod, `
		SELECT *
		FROM participation_methods
		WHERE participation_method_id = $1
	`, participationMethodId)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		return nil, er.ParticipationMethodNotFound
	default:
		return nil, err
	}
	return &participationMethod, nil
}
