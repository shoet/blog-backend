package create_blog

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/session"
)

type BlogRepository interface {
	Add(ctx context.Context, tx infrastructure.TX, blog *models.Blog) (models.BlogId, error)
	Get(ctx context.Context, tx infrastructure.TX, id models.BlogId) (*models.Blog, error)
	AddBlogTag(ctx context.Context, tx infrastructure.TX, blogId models.BlogId, tagId models.TagId) (int64, error)
	SelectTags(ctx context.Context, tx infrastructure.TX, tag string) ([]*models.Tag, error)
	AddTag(ctx context.Context, tx infrastructure.TX, tag string) (models.TagId, error)
}

type BlogService interface {
	Validate(ctx context.Context, userId models.UserId, blog *models.Blog) error
}

type Usecase struct {
	DB             infrastructure.DB
	BlogRepository BlogRepository
	BlogService    BlogService
}

func NewUsecase(
	db infrastructure.DB,
	blogRepository BlogRepository,
	blogService BlogService,
) *Usecase {
	return &Usecase{
		DB:             db,
		BlogRepository: blogRepository,
		BlogService:    blogService,
	}
}

func (u *Usecase) Run(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	sessionUserId, err := session.GetUserId(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to session.GetUserId: %w", err)
	}
	if err := u.BlogService.Validate(ctx, sessionUserId, blog); err != nil {
		return nil, fmt.Errorf("failed to BlogService.Validate: %w", err)
	}

	transactor := infrastructure.NewTransactionProvider(u.DB)

	result, err := transactor.DoInTx(ctx, func(tx infrastructure.TX) (interface{}, error) {
		// add tags
		var tagIds []models.TagId
		for _, tag := range blog.Tags {
			tags, err := u.BlogRepository.SelectTags(ctx, tx, tag)
			if err != nil {
				return nil, fmt.Errorf("failed to upsert tag: %w", err)
			}
			if len(tags) == 0 {
				tagId, err := u.BlogRepository.AddTag(ctx, tx, tag)
				if err != nil {
					return nil, fmt.Errorf("failed to add tag: %w", err)
				}
				tagIds = append(tagIds, tagId)
			} else {
				tagIds = append(tagIds, tags[0].Id)
			}
		}

		// add blog
		id, err := u.BlogRepository.Add(ctx, tx, blog)
		if err != nil {
			return nil, fmt.Errorf("failed to add blog: %w", err)
		}

		// add blogs_tags
		for _, tagId := range tagIds {
			_, err := u.BlogRepository.AddBlogTag(ctx, tx, id, tagId)
			if err != nil {
				return nil, fmt.Errorf("failed to add blogs_tags: %w", err)
			}
		}

		newBlog, err := u.BlogRepository.Get(ctx, tx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get blog: %w", err)
		}
		return newBlog, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}

	blog, ok := result.(*models.Blog)
	if !ok {
		return nil, fmt.Errorf("failed to type assertion: %w", err)
	}

	return blog, nil
}
