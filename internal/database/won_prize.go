package database

import (
	"context"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type WonPrizeRepository interface {
	CreateWonPrize(ctx context.Context, wonPrize database.CreateWonPrize) error
}

type wonPrizeRepository struct {
	db *sqlx.DB
}

func NewWonPrizeRepository(db *sqlx.DB) WonPrizeRepository {
	return &wonPrizeRepository{
		db: db,
	}
}

func (w *wonPrizeRepository) CreateWonPrize(ctx context.Context, wonPrize database.CreateWonPrize) error {
	tx, err := w.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO won_prize (prize_id, collaborator_id)
		VALUES (:prize_id, :collaborator_id)
	`, wonPrize)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE collaborator
		SET last_roll_time = current_timestamp
		WHERE collaborator_id = $1
	`, wonPrize.CollaboratorId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
