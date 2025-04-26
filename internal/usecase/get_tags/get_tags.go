package get_tags

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepository interface {
	ListTags(ctx context.Context, tx infrastructure.TX, option options.ListTagsOptions) ([]*models.Tag, error)
}

type Usecase struct {
	DB             infrastructure.DB
	BlogRepository BlogRepository
}

func NewUsecase(
	db infrastructure.DB,
	blogRepository BlogRepository,
) *Usecase {
	return &Usecase{
		DB:             db,
		BlogRepository: blogRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, option options.ListTagsOptions) (models.Tags, error) {
	transactor := infrastructure.NewTransactionProvider(u.DB)
	result, err := transactor.DoInTx(ctx, func(tx infrastructure.TX) (interface{}, error) {
		tags, err := u.BlogRepository.ListTags(ctx, tx, option)
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}
		return tags, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	tags, ok := result.([]*models.Tag)
	if !ok {
		return nil, fmt.Errorf("failed to assert result to models.Tags")
	}

	return tags, nil
}
