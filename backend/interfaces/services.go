package interfaces

import (
	"context"

	"github.com/shoet/blog/config"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogService interface {
	ListBlog(ctx context.Context, option options.ListBlogOptions) ([]*models.Blog, error)
	AddBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	DeleteBlog(ctx context.Context, id models.BlogId) error
	PutBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	Export(ctx context.Context) error
	GetBlog(ctx context.Context, id models.BlogId) (*models.Blog, error)
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	LoginAdmin(ctx context.Context, cfg *config.Config, email string, password string) (string, error)
	LoginSession(ctx context.Context, token string) (*models.User, error)
	// Signup(ctx context.Context, email string, password string) (string, error)
	// Signout(ctx context.Context, token string) error
	// Unsubscribe(ctx context.Context) error
}

type StorageService interface {
	GenerateThumbnailPutURL(fileName string) (string, string, error)
}
