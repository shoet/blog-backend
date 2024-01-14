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
		var usingTagsByOtherBlog models.BlogsTagsArray
		usingTagsByOtherBlog, err = u.BlogRepository.SelectBlogsTagsByOtherUsingBlog(ctx, tx, blog.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to select using tags: %w", err)
		}
		isUsing := func(tag string) bool { return slices.Contains(usingTagsByOtherBlog.TagNames(), tag) }

		var currentTags models.BlogsTagsArray
		currentTags, err := u.BlogRepository.SelectBlogsTags(ctx, tx, blog.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to select current tags: %w", err)
		}
		isCurrent := func(tag string) bool { return slices.Contains(currentTags.TagNames(), tag) }

		isContainsNew := func(tag string) bool { return slices.Contains(blog.Tags, tag) }

		for _, tag := range blog.Tags {
			if !isCurrent(tag) {
				tags, err := u.BlogRepository.SelectTags(ctx, tx, tag)
				if err != nil {
					return nil, fmt.Errorf("failed to select tag: %w", err)
				}
				// add tags
				var tagId models.TagId
				if len(tags) == 0 {
					tagId, err = u.BlogRepository.AddTag(ctx, tx, tag)
					if err != nil {
						return nil, fmt.Errorf("failed to add tag: %w", err)
					}
				} else {
					tagId = tags[0].Id
				}
				// add blogs_tags
				if _, err := u.BlogRepository.AddBlogTag(ctx, tx, blog.Id, tagId); err != nil {
					return nil, fmt.Errorf("failed to add blogs_tags: %w", err)
				}
			}
		}

		for _, tag := range currentTags {
			if isContainsNew(tag.Name) {
				continue
			}
			if !isUsing(tag.Name) {
				// delete tags
				if err := u.BlogRepository.DeleteTag(ctx, tx, tag.TagId); err != nil {
					return nil, fmt.Errorf("failed to delete tags: %w", err)
				}
				if err := u.BlogRepository.DeleteBlogsTags(ctx, tx, blog.Id, tag.TagId); err != nil {
					return nil, fmt.Errorf("failed to delete blogs_tags: %w", err)
				}
			}
		}

		// put blog
		id, err := u.BlogRepository.Put(ctx, tx, blog)
		if err != nil {
			return nil, fmt.Errorf("failed to put blog: %w", err)
		}

		newBlog, err := u.BlogRepository.Get(ctx, tx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get blog: %w", err)
		}

		return newBlog, nil
	})

	blog, ok := result.(*models.Blog)
	if !ok {
		return nil, fmt.Errorf("failed to type assertion: %w", err)
	}
	return blog, nil
}
