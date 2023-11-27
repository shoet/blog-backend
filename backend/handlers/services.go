package handlers

import (
	"context"

	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogManager interface {
	ListBlog(ctx context.Context, option options.ListBlogOptions) ([]*models.Blog, error)
	AddBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	DeleteBlog(ctx context.Context, id models.BlogId) error
	PutBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	Export(ctx context.Context) error
	GetBlog(ctx context.Context, id models.BlogId) (*models.Blog, error)
	ListTags(ctx context.Context, option options.ListTagsOptions) ([]*models.Tag, error)
}

type AuthManager interface {
	Login(ctx context.Context, email string, password string) (string, error)
	LoginSession(ctx context.Context, token string) (*models.User, error)
	// Signup(ctx context.Context, email string, password string) (string, error)
	// Signout(ctx context.Context, token string) error
	// Unsubscribe(ctx context.Context) error
}

type Storager interface {
	GenerateThumbnailPutURL(fileName string) (string, string, error)
	GenerateContentImagePutURL(fileName string) (string, string, error)
}
