package database

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserRepository interface {
	GetUsers(ctx context.Context, filter database.GetUsersFilter) ([]database.User, error)
	GetUser(ctx context.Context, userId string) (*database.User, error)
	GetUserFromFields(ctx context.Context, gameId string, fields database.UserFields) (*database.User, error)
	CreateUser(ctx context.Context, user database.CreateUser) (*database.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUsers(ctx context.Context, filter database.GetUsersFilter) ([]database.User, error) {
	query := sq.
		Select("*").
		From("users")

	if filter.GameId != nil {
		query = query.Where(sq.Eq{"game_id": *filter.GameId})
	}

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	var users []database.User
	err = r.db.SelectContext(ctx, &users, sqlQ, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUser(ctx context.Context, userId string) (*database.User, error) {
	var user database.User
	err := r.db.GetContext(ctx, &user, `
		SELECT *
		FROM users
		WHERE user_id = $1
	`, userId)
	switch {
	case err == nil:
		return &user, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, er.UserNotFound
	default:
		return nil, err
	}
}

func (r *userRepository) GetUserFromFields(ctx context.Context, gameId string, fields database.UserFields) (*database.User, error) {
	query := sq.
		Select("*").
		From("users").
		Where("game_id = ?", gameId)

	hasFilter := false
	if fields.Email != nil {
		query = query.Where(sq.Eq{"email": *fields.Email})
		hasFilter = true
	}
	if fields.Address != nil {
		query = query.Where(sq.Eq{"address": *fields.Address})
		hasFilter = true
	}
	if fields.Phone != nil {
		query = query.Where(sq.Eq{"phone": *fields.Phone})
		hasFilter = true
	}

	if !hasFilter {
		return nil, errors.New("no user fields provided")
	}

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate sql query: %w", err)
	}

	var user database.User
	err = r.db.GetContext(ctx, &user, sqlQ, args...)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		return nil, er.UserNotFound
	default:
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, create database.CreateUser) (*database.User, error) {
	query := sq.
		Select("*").
		From("users").
		Where(sq.Eq{"game_id": create.GameId})

	or := sq.Or{}
	if create.Email != nil {
		or = append(or, sq.Eq{"email": *create.Email})
	}
	if create.Address != nil {
		or = append(or, sq.Eq{"address": *create.Address})
	}
	if create.Phone != nil {
		or = append(or, sq.Eq{"phone": *create.Phone})
	}
	if len(or) > 0 {
		query = query.Where(or)
	}

	sqlQ, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate sql: %w", err)
	}

	var c database.User
	err = r.db.GetContext(ctx, &c, sqlQ, args...)
	switch {
	case err == nil:
		return nil, er.UserExists
	case errors.Is(err, sql.ErrNoRows):
		break
	default:
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	_, err = r.db.NamedExecContext(ctx, `
		INSERT INTO users (game_id, email, address, phone, additional_fields)
		VALUES (:game_id, :email, :address, :phone, :additional_fields)
	`, create)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	user, err := r.GetUserFromFields(ctx, create.GameId, create.UserFields)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}
