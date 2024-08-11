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
		Select("p.*").
		From("prizes p")

	if filter.GameId != nil {
		query = query.Where("game_id = ?", filter.GameId)
	}
	if filter.DrawMethodId != nil {
		query = query.
			LeftJoin("draw_methods_prizes USING (prize_id)").
			Where("draw_method_id = ?", filter.DrawMethodId)
	}
	if filter.AvailableOnly {
		subQuery := sq.
			Select("COUNT(*)").
			From("won_prizes wp").
			InnerJoin("draw_methods_prizes USING (prize_id)").
			Where("wp.prize_id = p.prize_id")

		if filter.DrawMethodId != nil {
			subQuery = subQuery.Where("draw_method_id = ?", filter.DrawMethodId)
		}

		sql, args, err := subQuery.ToSql()
		if err != nil {
			return nil, err
		}

		query = query.Where(fmt.Sprintf("p.count > (%s)", sql), args...)
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
