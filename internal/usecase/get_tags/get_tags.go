package get_tags

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepository interface {
	ListTags(ctx context.Context, tx infrastracture.TX, option options.ListTagsOptions) ([]*models.Tag, error)
}

type Usecase struct {
	DB             infrastracture.DB
	BlogRepository BlogRepository
}

func NewUsecase(
	db infrastracture.DB,
	blogRepository BlogRepository,
) *Usecase {
	return &Usecase{
		DB:             db,
		BlogRepository: blogRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, option options.ListTagsOptions) (models.Tags, error) {
	transactor := infrastracture.NewTransactionProvider(u.DB)
	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
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
