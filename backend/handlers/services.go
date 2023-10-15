package handlers

import (
	"context"

	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogService interface {
	ListBlog(ctx context.Context, option options.ListBlogOptions) ([]*models.Blog, error)
	GetBlog(ctx context.Context, id models.BlogId) (*models.Blog, error)
	AddBlog(ctx context.Context, blog *models.Blog) error
	DeleteBlog(ctx context.Context, id models.BlogId) error
	PutBlog(ctx context.Context, blog *models.Blog) error
}
