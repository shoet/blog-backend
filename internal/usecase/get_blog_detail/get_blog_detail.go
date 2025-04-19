package get_blog_detail

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type BlogRepository interface {
	Get(ctx context.Context, tx infrastracture.TX, id models.BlogId) (*models.Blog, error)
}

type CommentRepository interface {
	GetByBlogId(
		ctx context.Context, tx infrastracture.TX, blogId models.BlogId, excludeDeleted bool,
	) ([]*models.Comment, error)
}

type Usecase struct {
	DB                infrastracture.DB
	BlogRepository    BlogRepository
	CommentRepository CommentRepository
}

func NewUsecase(db infrastracture.DB, blogRepository BlogRepository, commentRepository CommentRepository) *Usecase {
	return &Usecase{
		DB:                db,
		BlogRepository:    blogRepository,
		CommentRepository: commentRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId) (*models.Blog, []*models.Comment, error) {
	blog, err := u.BlogRepository.Get(ctx, u.DB, blogId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get blog: %v", err)
	}
	if blog == nil {
		return nil, nil, nil
	}
	comments, err := u.CommentRepository.GetByBlogId(ctx, u.DB, blogId, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get comments: %v", err)
	}
	return blog, comments, nil

}
