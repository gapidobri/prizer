package database

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ParticipationMethodRepository interface {
	GetParticipationMethods(ctx context.Context, filter database.GetParticipationMethodsFilter) ([]database.ParticipationMethod, error)
	GetParticipationMethod(ctx context.Context, participationMethodId string) (*database.ParticipationMethod, error)
	UpdateParticipationMethod(ctx context.Context, participationMethodId string, participationMethod database.UpdateParticipationMethod) error
	LinkDrawMethod(ctx context.Context, participationMethodId string, drawMethodId string) error
	UnlinkDrawMethod(ctx context.Context, participationMethodId string, drawMethodId string) error
}

type participationMethodRepository struct {
	db *sqlx.DB
}

func NewParticipationMethodRepository(db *sqlx.DB) ParticipationMethodRepository {
	return &participationMethodRepository{
		db: db,
	}
}

func (r *participationMethodRepository) GetParticipationMethods(ctx context.Context, filter database.GetParticipationMethodsFilter) ([]database.ParticipationMethod, error) {
	query := sq.
		Select("*").
		From("participation_methods")

	if filter.GameId != nil {
		query = query.Where(sq.Eq{"game_id": filter.GameId})
	}

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	var participationMethods []database.ParticipationMethod
	err = r.db.SelectContext(ctx, &participationMethods, sqlQ, args...)
	if err != nil {
		return nil, err
	}
	return participationMethods, nil
}

func (r *participationMethodRepository) GetParticipationMethod(ctx context.Context, participationMethodId string) (*database.ParticipationMethod, error) {
	var participationMethod database.ParticipationMethod
	err := r.db.GetContext(ctx, &participationMethod, `
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

func (r *participationMethodRepository) UpdateParticipationMethod(ctx context.Context, participationMethodId string, participationMethod database.UpdateParticipationMethod) error {
	setMap := map[string]interface{}{
		"fields": &participationMethod.Fields,
	}
	err := mapstructure.Decode(participationMethod, &setMap)
	if err != nil {
		return err
	}

	query := sq.
		Update("participation_methods").
		Where(sq.Eq{"participation_method_id": participationMethodId}).
		SetMap(setMap)

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlQ, args...)
	return err
}

func (r *participationMethodRepository) LinkDrawMethod(ctx context.Context, participationMethodId string, drawMethodId string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO participation_methods_draw_methods (participation_method_id, draw_method_id)
		VALUES ($1, $2)
	`, participationMethodId, drawMethodId)
	return err
}

func (r *participationMethodRepository) UnlinkDrawMethod(ctx context.Context, participationMethodId string, drawMethodId string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM participation_methods_draw_methods
		WHERE participation_method_id = $1
			AND draw_method_id = $2
	`, participationMethodId, drawMethodId)
	return err
}
