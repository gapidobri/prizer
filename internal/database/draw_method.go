package database

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type DrawMethodRepository interface {
	GetDrawMethods(ctx context.Context, gameId string, filter database.GetDrawMethodsFilter) ([]database.DrawMethod, error)
}

type drawMethodRepository struct {
	db *sqlx.DB
}

func NewDrawMethodRepository(db *sqlx.DB) DrawMethodRepository {
	return &drawMethodRepository{
		db: db,
	}
}

func (d *drawMethodRepository) GetDrawMethods(ctx context.Context, gameId string, filter database.GetDrawMethodsFilter) ([]database.DrawMethod, error) {
	query := sq.
		Select("DISTINCT ON (draw_method_id) dm.*").
		From("collaboration_method").
		InnerJoin("collaboration_method_draw_method USING (collaboration_method_id)").
		InnerJoin("draw_method dm USING (draw_method_id)").
		Where("game_id = ?", gameId)

	if filter.CollaborationMethodId != nil {
		query = query.Where("collaboration_method_id = ?", filter.CollaborationMethodId)
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
