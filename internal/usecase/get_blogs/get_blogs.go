package get_blogs

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepository interface {
	List(ctx context.Context, tx infrastracture.TX, option *options.ListBlogOptions) ([]*models.Blog, error)
}

type Usecase struct {
	DB             infrastracture.DB
	BlogRepository BlogRepository
}

func NewUsecase(
	DB infrastracture.DB,
	blogRepository BlogRepository,
) *Usecase {
	return &Usecase{
		DB:             DB,
		BlogRepository: blogRepository,
	}
}

type GetBlogsInput struct {
	IsPublic *bool
	Tag      *string
	KeyWord  *string
}

func NewGetBlogsInput(
	isPublic *bool,
	tag *string,
	keyWord *string,
) *GetBlogsInput {
	input := new(GetBlogsInput)
	input.IsPublic = isPublic
	input.Tag = tag
	input.KeyWord = keyWord
	return input
}

func (u *Usecase) Run(ctx context.Context, input *GetBlogsInput) ([]*models.Blog, error) {

	transactor := infrastracture.NewTransactionProvider(u.DB)

	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		listOption := options.NewListBlogOptions(input.IsPublic, nil)
		var blogs models.Blogs
		blogs, err := u.BlogRepository.List(ctx, tx, listOption)
		if err != nil {
			return nil, fmt.Errorf("failed to list blogs: %v", err)
		}
		// タグの検索を優先する
		if input.Tag != nil {
			blogs = blogs.FilterByTag(*input.Tag)
		} else if input.KeyWord != nil {
			blogs = blogs.FilterByKeyword(*input.KeyWord)
		}
		return blogs.ToSlice(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get blogs: %v", err)
	}

	blogs, ok := result.([]*models.Blog)
	if !ok {
		return nil, fmt.Errorf("failed to cast []*models.Blog")
	}
	return blogs, nil
}
