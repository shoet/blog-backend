package handler

import (
	"context"

	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

//go:generate go run github.com/matryer/moq -out services_moq.go . BlogManager AuthManager Storager
type BlogManager interface {
	DeleteBlog(ctx context.Context, id models.BlogId) error
	PutBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	Export(ctx context.Context) error
	ListTags(ctx context.Context, option options.ListTagsOptions) ([]*models.Tag, error)
}

type AuthManager interface {
	Login(ctx context.Context, email string, password string) (string, error)
	LoginSession(ctx context.Context, token string) (*models.User, error)
}

type Storager interface {
	GenerateThumbnailPutURL(fileName string) (string, string, error)
	GenerateContentImagePutURL(fileName string) (string, string, error)
}
