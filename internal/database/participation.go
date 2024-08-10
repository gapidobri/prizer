package database

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type ParticipationRepository interface {
	CreateParticipation(ctx context.Context, createParticipation database.CreateParticipation) error
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

func (r *participationRepository) CreateParticipation(ctx context.Context, createParticipation database.CreateParticipation) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO participation (participation_method_id, user_id, fields)
		VALUES (:participation_method_id, :user_id, :fields)
	`, createParticipation)

	return err
}

func (r *participationRepository) GetParticipations(ctx context.Context, filter database.GetParticipationsFilter) ([]database.Participation, error) {
	query := sq.
		Select("*").
		From("participation")

	if filter.UserId != nil {
		query = query.Where(sq.Eq{"user_id": filter.UserId})
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
