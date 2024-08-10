package database

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type PrizeRepository interface {
	GetPrizes(ctx context.Context, filter database.GetPrizesFilter) ([]database.Prize, error)
}

type prizeRepository struct {
	db *sqlx.DB
}

func NewPrizeRepository(db *sqlx.DB) PrizeRepository {
	return &prizeRepository{
		db: db,
	}
}

func (p *prizeRepository) GetPrizes(ctx context.Context, filter database.GetPrizesFilter) ([]database.Prize, error) {
	query := sq.
		Select("*").
		From("prize")

	if filter.GameId != nil {
		query.Where("game_id = ?", filter.GameId)
	}
	if filter.AvailableOnly {
		subQuery, subArgs := sq.
			Select("COUNT(*)").
			From("won_prize wp").
			InnerJoin("prize p USING (prize_id)").
			PlaceholderFormat(sq.Dollar).
			MustSql()

		query = query.Where(fmt.Sprintf("count > (%s)", subQuery), subArgs...)
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var prizes []database.Prize
	err = p.db.SelectContext(ctx, &prizes, sql, args...)
	if err != nil {
		return nil, err
	}

	return prizes, nil
}
