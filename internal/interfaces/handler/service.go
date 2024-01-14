package handler

import (
	"context"

	"github.com/shoet/blog/internal/infrastracture/models"
)

//go:generate go run github.com/matryer/moq -out services_moq.go . BlogManager AuthManager Storager
type BlogManager interface {
	Export(ctx context.Context) error
}

type AuthManager interface {
	Login(ctx context.Context, email string, password string) (string, error)
	LoginSession(ctx context.Context, token string) (*models.User, error)
}

type Storager interface {
	GenerateThumbnailPutURL(fileName string) (string, string, error)
	GenerateContentImagePutURL(fileName string) (string, string, error)
}
