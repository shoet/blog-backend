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
	ListByTag(ctx context.Context, tx infrastracture.TX, tag string, isPublicOnly bool) (models.Blogs, error)
	ListByKeyword(ctx context.Context, tx infrastracture.TX, keyword string, isPublicOnly bool) (models.Blogs, error)
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
	Offset   *uint
	Limit    *uint
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

		if input.Tag != nil {
			// タグ検索
			b, err := u.BlogRepository.ListByTag(ctx, tx, *input.Tag, *input.IsPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by tag: %v", err)
			}
			blogs = b
		} else if input.KeyWord != nil {
			// キーワード検索
			b, err := u.BlogRepository.ListByKeyword(ctx, tx, *input.KeyWord, *input.IsPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by keyword: %v", err)
			}
			blogs = b
		} else {
			// 通常の検索
			b, err := u.BlogRepository.List(ctx, tx, listOption)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs: %v", err)
			}
			blogs = b
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
