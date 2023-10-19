package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db    *sqlx.DB
	user  UserRepository
	jwter JWTer
}

func NewAuthService(
	db *sqlx.DB, user UserRepository, jwter JWTer,
) (*AuthService, error) {
	return &AuthService{
		db:    db,
		user:  user,
		jwter: jwter,
	}, nil
}

func (a *AuthService) Login(
	ctx context.Context, email string, password string,
) (string, error) {

	// get user
	u, err := a.user.GetByEmail(ctx, a.db, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user by email: %w", err)
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("failed to compare password: %w", err)
	}

	// generate token and save session kvs
	token, err := a.jwter.GenerateToken(ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a *AuthService) LoginAdmin(
	ctx context.Context, cfg *config.Config, email string, password string,
) (string, error) {
	if cfg.AdminEmail != email {
		return "", fmt.Errorf("invalid admin email")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(cfg.AdminPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("failed to compare password: %w", err)
	}

	u, err := a.user.GetByEmail(ctx, a.db, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user by email: %w", err)
	}

	// generate token and save session kvs
	token, err := a.jwter.GenerateToken(ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, err
}

func (a *AuthService) LoginSession(
	ctx context.Context, token string,
) (*models.User, error) {
	return nil, nil
}
