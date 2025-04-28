package handler

import (
	"context"
	"net/http"

	"github.com/shoet/blog/internal/infrastructure/models"
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

type Cookier interface {
	SetCookie(w http.ResponseWriter, key string, value string) error
	ClearCookie(w http.ResponseWriter, key string)
}
