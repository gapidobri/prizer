package database

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type WonPrizeRepository interface {
	CreateWonPrize(ctx context.Context, wonPrize database.CreateWonPrize) error
	GetWonPrizes(ctx context.Context, filter database.GetWonPrizesFilter) ([]database.WonPrize, error)
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
	_, err := w.db.NamedExecContext(ctx, `
		INSERT INTO won_prizes (prize_id, participation_id)
		VALUES (:prize_id, :participation_id)
	`, wonPrize)
	return err
}

func (w *wonPrizeRepository) GetWonPrizes(ctx context.Context, filter database.GetWonPrizesFilter) ([]database.WonPrize, error) {
	query := sq.
		Select("*").
		From("won_prizes").
		InnerJoin("prizes p USING (prize_id)").
		InnerJoin("participations USING (participation_id)").
		InnerJoin("users USING (user_id)").
		OrderBy("created_at DESC")

	if filter.PrizeId != nil {
		query = query.Where(sq.Eq{"prize_id": filter.PrizeId})
	}
	if filter.UserId != nil {
		query = query.Where(sq.Eq{"user_id": filter.UserId})
	}
	if filter.GameId != nil {
		query = query.Where(sq.Eq{"p.game_id": filter.GameId})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var prizes []database.WonPrize
	err = w.db.SelectContext(ctx, &prizes, sql, args...)
	if err != nil {
		return nil, err
	}
	return prizes, nil
}
