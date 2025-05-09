package update_public_status

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type BlogRepository interface {
	UpdatePublicStatus(
		ctx context.Context, tx infrastructure.TX, blogId models.BlogId, isPublic bool,
	) (*models.Blog, error)
}

type Usecase struct {
	DB             infrastructure.DB
	blogRepository BlogRepository
}

func NewUsecase(
	db infrastructure.DB, blogRepository BlogRepository) *Usecase {
	return &Usecase{
		DB:             db,
		blogRepository: blogRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId, isPublic bool) (*models.Blog, error) {
	blog, err := u.blogRepository.UpdatePublicStatus(ctx, u.DB, blogId, isPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog public status: %w", err)
	}
	return blog, nil
}
