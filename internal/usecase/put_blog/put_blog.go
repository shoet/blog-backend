package put_blog

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/session"
	"golang.org/x/exp/slices"
)

type BlogRepository interface {
	SelectBlogsTags(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.BlogsTags, error)
	SelectBlogsTagsByOtherUsingBlog(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.BlogsTags, error)
	SelectTags(ctx context.Context, tx infrastracture.TX, tag string) ([]*models.Tag, error)
	AddTag(ctx context.Context, tx infrastracture.TX, tag string) (models.TagId, error)
	AddBlogTag(ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId) (int64, error)
	DeleteTag(ctx context.Context, tx infrastracture.TX, tagId models.TagId) error
	DeleteBlogsTags(ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId) error
	Put(ctx context.Context, tx infrastracture.TX, blog *models.Blog) (models.BlogId, error)
	Get(ctx context.Context, tx infrastracture.TX, id models.BlogId) (*models.Blog, error)
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

func (u *Usecase) Run(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	sessionUserId, err := session.GetUserId(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to session.GetUserId: %w", err)
	}
	if sessionUserId != blog.AuthorId {
		return nil, fmt.Errorf("can't update other user's blog")
	}

	transactor := infrastracture.NewTransactionProvider(u.DB)
	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		// このブログに紐づいているタグで、他のブログで使用されているタグを取得する
		var usingTagsByOtherBlog models.BlogsTagsArray
		usingTagsByOtherBlog, err = u.BlogRepository.SelectBlogsTagsByOtherUsingBlog(ctx, tx, blog.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to select using tags: %w", err)
		}

		// このブログに紐づいているタグを取得する
		var currentTags models.BlogsTagsArray
		currentTags, err := u.BlogRepository.SelectBlogsTags(ctx, tx, blog.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to select current tags: %w", err)
		}

		// 新規のタグ追加
		for _, tag := range blog.Tags {
			if !currentTags.Contains(tag) {
				tags, err := u.BlogRepository.SelectTags(ctx, tx, tag)
				if err != nil {
					return nil, fmt.Errorf("failed to select tag: %w", err)
				}
				// タグを追加
				var tagId models.TagId
				if len(tags) == 0 {
					tagId, err = u.BlogRepository.AddTag(ctx, tx, tag)
					if err != nil {
						return nil, fmt.Errorf("failed to add tag: %w", err)
					}
				} else {
					tagId = tags[0].Id
				}
				// タグのリレーションを追加
				if _, err := u.BlogRepository.AddBlogTag(ctx, tx, blog.Id, tagId); err != nil {
					return nil, fmt.Errorf("failed to add blogs_tags: %w", err)
				}
			}
		}

		// 不要となったタグの削除
		for _, tag := range currentTags {
			if slices.Contains(blog.Tags, tag.Name) {
				continue
			}
			if !usingTagsByOtherBlog.Contains(tag.Name) {
				// 他のブログで使用されていないタグは削除
				if err := u.BlogRepository.DeleteTag(ctx, tx, tag.TagId); err != nil {
					return nil, fmt.Errorf("failed to delete tags: %w", err)
				}
				if err := u.BlogRepository.DeleteBlogsTags(ctx, tx, blog.Id, tag.TagId); err != nil {
					return nil, fmt.Errorf("failed to delete blogs_tags: %w", err)
				}
			}
		}

		// ブログの更新
		id, err := u.BlogRepository.Put(ctx, tx, blog)
		if err != nil {
			return nil, fmt.Errorf("failed to put blog: %w", err)
		}

		// 更新後のブログを取得
		newBlog, err := u.BlogRepository.Get(ctx, tx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get blog: %w", err)
		}

		return newBlog, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}

	blog, ok := result.(*models.Blog)
	if !ok {
		return nil, fmt.Errorf("failed to type assertion: %w", err)
	}
	return blog, nil
}
