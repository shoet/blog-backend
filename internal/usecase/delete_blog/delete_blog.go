package delete_blog

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/session"
	"golang.org/x/exp/slices"
)

type BlogRepository interface {
	Get(ctx context.Context, tx infrastracture.TX, id models.BlogId) (*models.Blog, error)
	Delete(ctx context.Context, tx infrastracture.TX, id models.BlogId) error
	SelectBlogsTagsByOtherUsingBlog(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.BlogsTags, error)
	SelectBlogsTags(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.BlogsTags, error)
	DeleteTag(ctx context.Context, tx infrastracture.TX, tagId models.TagId) error
	DeleteBlogsTags(ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId) error
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

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId) (models.BlogId, error) {
	sessionUserId, err := session.GetUserId(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to session.GetUserId: %w", err)
	}

	transactor := infrastracture.NewTransactionProvider(u.DB)

	result, err := transactor.DoInTx(ctx, func(tx infrastracture.TX) (interface{}, error) {
		blog, err := u.BlogRepository.Get(ctx, tx, blogId)
		if err != nil {
			return 0, fmt.Errorf("failed to BlogRepository.Get: %w", err)
		}

		if blog.AuthorId != sessionUserId {
			return 0, fmt.Errorf("can't delete other user's blog")
		}

		// delete blogs_tags -----------------
		// select using other blog tags
		var usingTags models.BlogsTagsArray
		usingTags, err = u.BlogRepository.SelectBlogsTagsByOtherUsingBlog(ctx, tx, blog.Id)
		if err != nil {
			return 0, fmt.Errorf("failed to select using tags: %w", err)
		}

		//  select will delete tags
		blogsTags, err := u.BlogRepository.SelectBlogsTags(ctx, tx, blog.Id)
		if err != nil {
			return 0, fmt.Errorf("failed to select blogs_tags: %w", err)
		}
		var willDeleteTags []models.TagId
		for _, t := range blogsTags {
			if !slices.Contains(usingTags.TagIds(), t.TagId) {
				willDeleteTags = append(willDeleteTags, t.TagId)
			}
		}

		for _, tagId := range willDeleteTags {
			// delete tags
			if err := u.BlogRepository.DeleteTag(ctx, tx, tagId); err != nil {
				return 0, fmt.Errorf("failed to delete tags: %w", err)
			}
			// delete blogs_tags
			if err := u.BlogRepository.DeleteBlogsTags(ctx, tx, blog.Id, tagId); err != nil {
				return 0, fmt.Errorf("failed to delete blogs_tags: %w", err)
			}
		}

		// delete blogs ----------------------
		err = u.BlogRepository.Delete(ctx, tx, blog.Id)
		if err != nil {
			return 0, fmt.Errorf("failed to delete blog: %w", err)
		}

		return blog.Id, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to delete blog: %w", err)
	}

	blogId, ok := result.(models.BlogId)
	if !ok {
		return 0, fmt.Errorf("failed to type assertion: %w", err)
	}
	return blogId, nil

}
