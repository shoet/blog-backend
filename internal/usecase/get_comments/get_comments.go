package get_comments

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type CommmentRepository interface {
	GetByBlogId(
		ctx context.Context,
		tx infrastructure.TX,
		blogId models.BlogId,
		excludeDeleted bool,
	) ([]*models.Comment, error)
}

type Usecase struct {
	DB                infrastructure.DB
	commentRepository CommmentRepository
}

func NewUsecase(
	db infrastructure.DB,
	commentRepository CommmentRepository,
) *Usecase {
	return &Usecase{
		DB:                db,
		commentRepository: commentRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId) ([]*models.Comment, error) {
	comments, err := u.commentRepository.GetByBlogId(ctx, u.DB, blogId, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	return comments, nil
}
