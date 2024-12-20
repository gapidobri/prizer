package database

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PrizeRepository interface {
	GetPrizes(ctx context.Context, filter database.GetPrizesFilter) ([]database.Prize, error)
	CreatePrize(ctx context.Context, prize database.CreatePrize) error
	UpdatePrize(ctx context.Context, prizeId string, prize database.UpdatePrize) error
	DeletePrize(ctx context.Context, prizeId string) error
}

type prizeRepository struct {
	db *sqlx.DB
}

func NewPrizeRepository(db *sqlx.DB) PrizeRepository {
	return &prizeRepository{
		db: db,
	}
}

func (r *prizeRepository) GetPrizes(ctx context.Context, filter database.GetPrizesFilter) ([]database.Prize, error) {
	query := sq.
		Select("p.*, COUNT(wp.prize_id) AS won_count").
		From("prizes p").
		LeftJoin("won_prizes wp USING (prize_id)").
		GroupBy("p.prize_id")

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

			sqlQ, args, err := subQuery.ToSql()
			if err != nil {
				return nil, err
			}

			query = query.Where(fmt.Sprintf("p.count > (%s)", sqlQ), args...)

			if filter.UserId != nil {
				userPrizeLimitQuery := sq.
					Select("prize_id").
					From("won_prizes").
					InnerJoin("participations USING (participation_id)").
					InnerJoin("draw_methods_prizes USING (prize_id)").
					Where(sq.Eq{
						"user_id":        filter.UserId,
						"draw_method_id": filter.DrawMethodId,
					}).
					GroupBy("prize_id", "user_prize_limit").
					Having("COUNT(prize_id) >= user_prize_limit")

				sqlQ, args, err = userPrizeLimitQuery.ToSql()
				if err != nil {
					return nil, err
				}

				query = query.Where(fmt.Sprintf("p.prize_id NOT IN (%s)", sqlQ), args...)
			}
		}

	}

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var prizes []database.Prize
	err = r.db.SelectContext(ctx, &prizes, sqlQ, args...)
	if err != nil {
		return nil, err
	}

	return prizes, nil
}

func (r *prizeRepository) CreatePrize(ctx context.Context, prize database.CreatePrize) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO prizes (game_id, name, description, image_url, count)
		VALUES (:game_id, :name, :description, :image_url, :count)
	`, prize)
	return err
}

func (r *prizeRepository) UpdatePrize(ctx context.Context, prizeId string, prize database.UpdatePrize) error {
	var setMap map[string]interface{}
	err := mapstructure.Decode(prize, &setMap)
	if err != nil {
		return err
	}

	query := sq.
		Update("prizes").
		Where(sq.Eq{"prize_id": prizeId}).
		SetMap(setMap)

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlQ, args...)
	return err
}

func (r *prizeRepository) DeletePrize(ctx context.Context, prizeId string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM prizes
	       WHERE prize_id = $1
	`, prizeId)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return er.PrizeNotFound
	default:
		return err
	}
}
