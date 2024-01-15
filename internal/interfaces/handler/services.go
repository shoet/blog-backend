package handler

import (
	"context"

	"github.com/shoet/blog/internal/infrastracture/models"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	LoginSession(ctx context.Context, token string) (*models.User, error)
}

type JWTService interface {
	GenerateToken(ctx context.Context, u *models.User) (string, error)
	VerifyToken(ctx context.Context, token string) (models.UserId, error)
}

type ContentsService interface {
	GenerateThumbnailPutURL(fileName string) (presignedUrl, objectUrl string, err error)
	GenerateContentImagePutURL(fileName string) (presignedUrl, objectUrl string, err error)
}
