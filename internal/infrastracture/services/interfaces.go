package services

import (
	"context"

	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
)

//go:generate go run github.com/matryer/moq -out service_moq.go . JWTer KVSer

type UserRepository interface {
	Add(ctx context.Context, db repository.Execer, user *models.User) (*models.User, error)
	Get(ctx context.Context, db repository.Queryer, id models.UserId) (*models.User, error)
	GetByEmail(ctx context.Context, db repository.Queryer, email string) (*models.User, error)
	// Delete(ctx context.Context, db store.Execer, id models.UserId) error
	// Put(ctx context.Context, db store.Execer, user *models.User) error
}

type JWTer interface {
	GenerateToken(ctx context.Context, u *models.User) (string, error)
	VerifyToken(ctx context.Context, token string) (models.UserId, error)
}

type KVSer interface {
	Save(ctx context.Context, key string, value string) error
	Load(ctx context.Context, key string) (string, error)
}
