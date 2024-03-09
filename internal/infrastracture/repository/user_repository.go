package repository

import (
	"context"
	"fmt"

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
	sql := `
	SELECT
		id, name, created, modified
	FROM users
	WHERE id = $1
	;
	`
	var users []*models.User
	if err := tx.SelectContext(ctx, &users, sql, id); err != nil {
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
	sql := `
	SELECT
		id, name, email, password, created, modified
	FROM users
	WHERE email = $1
	;
	`
	var users []*models.User
	if err := tx.SelectContext(ctx, &users, sql, email); err != nil {
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
	sql := `
	INSERT INTO users
		(name, email, password)
	VALUES
		($1, $2, $3)
	RETURNING id
	;
	`

	row := tx.QueryRowxContext(
		ctx,
		sql,
		user.Name, user.Email, user.Password)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to insert user: %w", row.Err())
	}
	var userId models.UserId
	err := row.Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	user.Id = userId
	return user, nil
}
