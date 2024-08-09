package database

import (
	"context"
	"database/sql"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type CollaborationMethodRepository interface {
	GetCollaborationMethod(ctx context.Context, collaborationMethodId string) (*database.CollaborationMethod, error)
}

type collaborationMethodRepository struct {
	db *sqlx.DB
}

func NewCollaborationMethodRepository(db *sqlx.DB) CollaborationMethodRepository {
	return &collaborationMethodRepository{
		db: db,
	}
}

func (c *collaborationMethodRepository) GetCollaborationMethod(ctx context.Context, collaborationMethodId string) (*database.CollaborationMethod, error) {
	var collaborationMethod database.CollaborationMethod
	err := c.db.GetContext(ctx, &collaborationMethod, `
		SELECT *
		FROM collaboration_method
		WHERE collaboration_method_id = $1
	`, collaborationMethodId)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		return nil, er.CollaborationMethodNotFount
	default:
		return nil, err
	}
	return &collaborationMethod, nil
}
