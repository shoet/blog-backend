package get_blogs

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepository interface {
	List(ctx context.Context, db infrastracture.TX, option options.ListBlogOptions) ([]*models.Blog, error)
}

type Usecase struct {
	DB             infrastracture.DB
	BlogRepository BlogRepository
}

func NewUsecase(
	blogRepository BlogRepository,
) *Usecase {
	return &Usecase{
		BlogRepository: blogRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, isPublic bool) (*models.Blogs, error) {

	transactor := infrastracture.NewTransactionProvider(u.DB)

	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		listOption := options.ListBlogOptions{
			IsPublic: isPublic,
		}
		blogs, err := u.BlogRepository.List(ctx, tx, listOption)
		if err != nil {
			return nil, err
		}
		return blogs, nil
	})
	if err != nil {
		return nil, err
	}

	blogs, ok := result.(*models.Blogs)
	if !ok {
		return nil, fmt.Errorf("failed to cast []*models.Blog")
	}
	return blogs, nil
}
