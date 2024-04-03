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
	IsPublicOnly *bool
	Tag          *string
	KeyWord      *string
	OffsetBlogId *models.BlogId
	Limit        *uint
}

func NewGetBlogsInput(
	isPublicOnly *bool,
	tag *string,
	keyWord *string,
	offsetBlogId *models.BlogId,
	limit *uint,
) *GetBlogsInput {
	input := new(GetBlogsInput)
	input.IsPublicOnly = isPublicOnly
	input.Tag = tag
	input.KeyWord = keyWord
	input.OffsetBlogId = offsetBlogId
	input.Limit = limit
	return input
}

func (u *Usecase) Run(ctx context.Context, input *GetBlogsInput) ([]*models.Blog, error) {

	transactor := infrastracture.NewTransactionProvider(u.DB)

	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		listOption, err := options.NewListBlogOptions(input.IsPublicOnly, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create list option: %v", err)
		}
		var blogs models.Blogs

		if input.Tag != nil {
			// タグ検索
			b, err := u.BlogRepository.ListByTag(ctx, tx, *input.Tag, *input.IsPublicOnly)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by tag: %v", err)
			}
			blogs = b
		} else if input.KeyWord != nil {
			// キーワード検索
			b, err := u.BlogRepository.ListByKeyword(ctx, tx, *input.KeyWord, *input.IsPublicOnly)
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
