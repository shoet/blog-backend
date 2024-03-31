package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type UserRepository struct {
	Clocker clocker.Clocker
}

func NewUserRepository(clocker clocker.Clocker) (*UserRepository, error) {
	return &UserRepository{Clocker: clocker}, nil
}

func (t *UserRepository) Get(
	ctx context.Context, tx infrastracture.TX, id models.UserId,
) (*models.User, error) {
	builder := GetQueryBuilderPostgres()
	users_builder := builder.
		Select("id, name, created, modified").
		From("users").
		Where(sq.Eq{"id": id})
	sql, args, err := users_builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var users []*models.User
	if err := tx.SelectContext(ctx, &users, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return users[0], nil
}

var ErrUserNotFound = fmt.Errorf("user not found")

func (u *UserRepository) GetByEmail(
	ctx context.Context, tx infrastracture.TX, email string,
) (*models.User, error) {
	builder := GetQueryBuilderPostgres()
	sql, args, err := builder.
		Select("id, name, email, password, created, modified").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var users []*models.User
	if err := tx.SelectContext(ctx, &users, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

func (u *UserRepository) Add(
	ctx context.Context, tx infrastracture.TX, user *models.User,
) (*models.User, error) {
	builder := GetQueryBuilderPostgres()
	sql, args, err := builder.
		Insert("users").
		Columns("name", "email", "password").
		Values(user.Name, user.Email, user.Password).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	row := tx.QueryRowxContext(
		ctx,
		sql,
		args...)
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
