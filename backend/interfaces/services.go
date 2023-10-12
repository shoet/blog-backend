package interfaces

import (
	"context"

	"github.com/shoet/blog/models"
)

type ListBlogOptions struct {
	AuthorId models.UserId
	Tags     []models.TagId
	IsPublic bool
}

type BlogService interface {
	ListBlog(ctx context.Context, option *ListBlogOptions) ([]*models.Blog, error)
	AddBlog(ctx context.Context, blog *models.Blog) error
	DeleteBlog(ctx context.Context, id models.BlogId) error
	PutBlog(ctx context.Context, blog *models.Blog) error
	Export(ctx context.Context) (error)
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Signup(ctx context.Context, email string, password string) (string, error)
	Signout(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context) error
}
