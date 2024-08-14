package database

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type DrawMethodRepository interface {
	GetDrawMethods(ctx context.Context, filter database.GetDrawMethodsFilter) ([]database.DrawMethod, error)
}

type drawMethodRepository struct {
	db *sqlx.DB
}

func NewDrawMethodRepository(db *sqlx.DB) DrawMethodRepository {
	return &drawMethodRepository{
		db: db,
	}
}

func (d *drawMethodRepository) GetDrawMethods(ctx context.Context, filter database.GetDrawMethodsFilter) ([]database.DrawMethod, error) {
	query := sq.
		Select("dm.*").
		From("draw_methods dm")

	if filter.GameId != nil {
		query = query.Where(sq.Eq{"dm.game_id": filter.GameId})
	}
	if filter.ParticipationMethodId != nil {
		query = query.
			InnerJoin("participation_methods_draw_methods USING (draw_method_id)").
			Where(sq.Eq{"participation_method_id": filter.ParticipationMethodId})
	}

	sql, args := query.
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var drawMethods []database.DrawMethod
	err := d.db.SelectContext(ctx, &drawMethods, sql, args...)
	if err != nil {
		return nil, err
	}
	return drawMethods, nil
}
