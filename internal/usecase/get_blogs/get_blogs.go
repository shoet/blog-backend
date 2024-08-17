package get_blogs

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepository interface {
	List(
		ctx context.Context, tx infrastracture.TX, option *options.ListBlogOptions,
	) ([]*models.Blog, error)

	ListByTag(
		ctx context.Context, tx infrastracture.TX, tag string, option *options.ListBlogOptions,
	) (models.Blogs, error)

	ListByKeyword(
		ctx context.Context, tx infrastracture.TX, keyword string, option *options.ListBlogOptions,
	) (models.Blogs, error)
}

// get_blogs.Usecaseはブログ一覧を取得するユースケースです。
// ページングはカーソル方式で実装しています。
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
	Tag           *string
	KeyWord       *string
	IsPublicOnly  *bool
	CursorId      *models.BlogId
	PageDirection *string
	Limit         *int64
}

func (u *Usecase) Run(
	ctx context.Context, input *GetBlogsInput,
) (blogs []*models.Blog, prevEOF bool, nextEOF bool, err error) {

	transactor := infrastracture.NewTransactionProvider(u.DB)

	option, err := options.NewListBlogOptions(input.IsPublicOnly, input.CursorId, input.Limit, input.PageDirection)
	if err != nil {
		return nil, false, false, fmt.Errorf("failed to create list option: %v", err)
	}
	// 次のページが存在するか判定するためにLimit+1で取得する
	option.Limit++

	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		var blogs models.Blogs

		if input.Tag != nil {
			// タグ検索
			b, err := u.BlogRepository.ListByTag(ctx, tx, *input.Tag, option)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by tag: %v", err)
			}
			blogs = b
		} else if input.KeyWord != nil {
			// キーワード検索
			b, err := u.BlogRepository.ListByKeyword(ctx, tx, *input.KeyWord, option)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs by keyword: %v", err)
			}
			blogs = b
		} else {
			// 通常の検索
			b, err := u.BlogRepository.List(ctx, tx, option)
			if err != nil {
				return nil, fmt.Errorf("failed to list blogs: %v", err)
			}
			blogs = b
		}

		return blogs.ToSlice(), nil
	})

	if err != nil {
		return nil, false, false, fmt.Errorf("failed to get blogs: %v", err)
	}

	blogs, ok := result.([]*models.Blog)
	if !ok {
		return nil, false, false, fmt.Errorf("failed to cast []*models.Blog")
	}

	var isEOF = false
	if len(blogs) <= int(option.Limit-1) {
		// Limit+1で取得しているため、Limit以下の場合はEOF
		isEOF = true
	} else {
		// Limit+1で取得しているため、Limitを超える場合は最後の要素を削除
		if option.PageDirection == "prev" {
			blogs = blogs[1:]
		} else {
			blogs = blogs[:len(blogs)-1]
		}
	}
	return blogs,
		option.PageDirection == "prev" && isEOF,
		option.PageDirection == "next" && isEOF,
		nil
}
