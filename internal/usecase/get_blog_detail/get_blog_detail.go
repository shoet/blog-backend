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
	GetByBlogId(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.Comment, error)
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

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId) (*models.Blog, error) {
	transactor := infrastracture.NewTransactionProvider(u.DB)

	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (any, error) {
		blog, err := u.BlogRepository.Get(ctx, tx, blogId)
		if err != nil {
			return nil, fmt.Errorf("failed to get blog: %v", err)
		}
		comments, err := u.CommentRepository.GetByBlogId(ctx, tx, blogId)
		if err != nil {
			return nil, fmt.Errorf("failed to get comments: %v", err)
		}
		blog.Comments = comments
		return blog, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %v", err)
	}
	blog, ok := result.(*models.Blog)
	if !ok {
		return nil, fmt.Errorf("failed to cast *models.Blog")
	}
	return blog, nil

}
