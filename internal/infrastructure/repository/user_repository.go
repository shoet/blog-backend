package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type UserRepository struct {
	Clocker clocker.Clocker
}

func NewUserRepository(clocker clocker.Clocker) (*UserRepository, error) {
	return &UserRepository{Clocker: clocker}, nil
}

func (t *UserRepository) Get(
	ctx context.Context, tx infrastructure.TX, id models.UserId,
) (*models.User, error) {
	sql, params, err := goqu.
		From("users").
		Select(
			"id", "name", "created", "modified",
		).
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var users []*models.User
	if err := tx.SelectContext(ctx, &users, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return users[0], nil
}

var ErrUserNotFound = fmt.Errorf("user not found")

func (u *UserRepository) GetByEmail(
	ctx context.Context, tx infrastructure.TX, email string,
) (*models.User, error) {
	sql, params, err := goqu.
		From("users").
		Select(
			"id", "name", "email", "password", "created", "modified",
		).
		Where(goqu.Ex{"email": email}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var users []*models.User
	if err := tx.SelectContext(ctx, &users, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

func (u *UserRepository) Add(
	ctx context.Context, tx infrastructure.TX, user *models.User,
) (*models.User, error) {
	sql, params, err := goqu.
		Insert("users").
		Cols("name", "email", "password").
		Vals(goqu.Vals{user.Name, user.Email, user.Password}).
		Returning("id").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	row := tx.QueryRowxContext(ctx, sql, params...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to insert user: %w", row.Err())
	}
	var userId models.UserId
	if err := row.Scan(&userId); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	user.Id = userId
	return user, nil
}
