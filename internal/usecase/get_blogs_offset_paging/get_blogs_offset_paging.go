package get_blogs_offset_paging

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepositoryOffset interface {
	List(
		ctx context.Context, tx infrastracture.TX, option *options.ListBlogOptions,
	) (models.Blogs, error)

	ListByTag(
		ctx context.Context, tx infrastracture.TX, tag string, option *options.ListBlogOptions,
	) (models.Blogs, error)

	ListByKeyword(
		ctx context.Context, tx infrastracture.TX, keyword string, option *options.ListBlogOptions,
	) (models.Blogs, error)
}

// get_blogs_offset_paging.Usecaseはブログ一覧を取得するユースケースです。
// ページングはオフセット方式で実装しています。
type Usecase struct {
	DB                   infrastracture.DB
	BlogRepositoryOffset BlogRepositoryOffset
}

func NewUsecase(
	DB infrastracture.DB,
	blogRepositoryOffset BlogRepositoryOffset,
) *Usecase {
	return &Usecase{
		DB:                   DB,
		BlogRepositoryOffset: blogRepositoryOffset,
	}
}

type Input struct {
	Tag          *string
	KeyWord      *string
	IsPublicOnly *bool
	Limit        *int64
	Page         *int64
}

func (u *Usecase) Run(ctx context.Context, input *Input) ([]*models.Blog, error) {
	transactor := infrastracture.NewTransactionProvider(u.DB)

	option, err := options.NewListBlogOffsetOptions(input.IsPublicOnly, input.Limit, input.Page)
	if err != nil {
		return nil, fmt.Errorf("failed to create list option: %v", err)
	}
	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		var blogs models.Blogs

		if input.Tag != nil {
			// タグ検索
			b, err := u.BlogRepositoryOffset.ListByTag(ctx, tx, *input.Tag, option)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by tag: %v", err)
			}
			blogs = b
		} else if input.KeyWord != nil {
			// キーワード検索
			b, err := u.BlogRepositoryOffset.ListByKeyword(ctx, tx, *input.KeyWord, option)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by keyword: %v", err)
			}
			blogs = b
		} else {
			// 通常の検索
			b, err := u.BlogRepositoryOffset.List(ctx, tx, option)
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
