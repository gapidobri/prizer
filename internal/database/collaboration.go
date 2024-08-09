package database

import (
	"context"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type CollaborationRepository interface{}

type collaborationRepository struct {
	db *sqlx.DB
}

func NewCollaborationRepository(db *sqlx.DB) CollaborationRepository {
	return &collaborationRepository{
		db: db,
	}
}

func (c *collaborationRepository) CreateCollaboration(ctx context.Context, createCollaboration database.CreateCollaboration) error {
	_, err := c.db.NamedExecContext(ctx, `
		INSERT INTO collaboration (collaboration_method_id, collaborator_id)
		VALUES (:collaboration_method_id, :collaborator_id)
	`, createCollaboration)

	return err
}
