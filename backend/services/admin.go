package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/models"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	db   *sqlx.DB
	user UserRepository
}

func NewAdminService(
	db *sqlx.DB, user UserRepository,
) (*AdminService, error) {
	return &AdminService{
		db:   db,
		user: user,
	}, nil
}

func (a *AdminService) SeedAdminUser(
	ctx context.Context, cfg *config.Config,
) (*models.User, error) {
	// hash password
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &models.User{
		Email:    cfg.AdminEmail,
		Password: string(passwordHashed),
		Name:     cfg.AdminName,
	}
	u, err := a.user.Add(ctx, a.db, user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %w", err)
	}
	// masking data
	u.Password = ""
	u.Id = 0
	return u, nil
}
