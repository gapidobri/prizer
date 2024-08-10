package database

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ParticipationRepository interface {
	CreateParticipation(ctx context.Context, createParticipation database.CreateParticipation) (*database.Participation, error)
	GetParticipations(ctx context.Context, filter database.GetParticipationsFilter) ([]database.Participation, error)
}

type participationRepository struct {
	db *sqlx.DB
}

func NewParticipationRepository(db *sqlx.DB) ParticipationRepository {
	return &participationRepository{
		db: db,
	}
}

func (r *participationRepository) CreateParticipation(ctx context.Context, createParticipation database.CreateParticipation) (*database.Participation, error) {
	res, err := r.db.NamedQueryContext(ctx, `
		INSERT INTO participations (participation_method_id, user_id, fields)
		VALUES (:participation_method_id, :user_id, :fields)
		RETURNING *
	`, createParticipation)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, errors.New("participation not found after creating")
	}

	var participation database.Participation
	err = res.StructScan(&participation)
	if err != nil {
		return nil, err
	}

	return &participation, nil
}

func (r *participationRepository) GetParticipations(ctx context.Context, filter database.GetParticipationsFilter) ([]database.Participation, error) {
	query := sq.
		Select("*").
		From("participations")

	if filter.UserId != nil {
		query = query.Where(sq.Eq{"user_id": filter.UserId})
	}
	if filter.ParticipationMethodId != nil {
		query = query.Where(sq.Eq{"participation_method_id": filter.ParticipationMethodId})
	}
	if filter.From != nil {
		query = query.Where(sq.GtOrEq{"created_at": filter.From})
	}
	if filter.To != nil {
		query = query.Where(sq.LtOrEq{"created_at": filter.To})
	}
	if filter.Fields != nil {
		for key, value := range *filter.Fields {
			query = query.Where(sq.Eq{fmt.Sprintf("fields->>'%s'", key): value})
		}
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	var participations []database.Participation
	err = r.db.SelectContext(ctx, &participations, sql, args...)
	if err != nil {
		return nil, err
	}
	return participations, err
}
