package get_blog_detail

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type BlogRepository interface {
	Get(ctx context.Context, tx infrastructure.TX, id models.BlogId) (*models.Blog, error)
}

type CommentRepository interface {
	GetByBlogId(
		ctx context.Context, tx infrastructure.TX, blogId models.BlogId, excludeDeleted bool,
	) ([]*models.Comment, error)
}

type Usecase struct {
	DB                infrastructure.DB
	BlogRepository    BlogRepository
	CommentRepository CommentRepository
}

func NewUsecase(db infrastructure.DB, blogRepository BlogRepository, commentRepository CommentRepository) *Usecase {
	return &Usecase{
		DB:                db,
		BlogRepository:    blogRepository,
		CommentRepository: commentRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId) (*models.Blog, error) {
	blog, err := u.BlogRepository.Get(ctx, u.DB, blogId)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %v", err)
	}
	if blog == nil {
		return nil, nil
	}
	return blog, nil

}
