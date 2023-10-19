package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/interfaces"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

func NewBlogService(db *sqlx.DB, blog interfaces.BlogRepository) *BlogService {
	return &BlogService{
		db:   db,
		blog: blog,
	}
}

type BlogService struct {
	db   *sqlx.DB
	blog interfaces.BlogRepository
}

func (b *BlogService) AddBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// add tags
	var tagIds []models.TagId
	for _, tag := range blog.Tags {
		tags, err := b.blog.SelectTags(ctx, tx, tag)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert tag: %w", err)
		}
		if len(tags) == 0 {
			tagId, err := b.blog.AddTag(ctx, tx, tag)
			if err != nil {
				return nil, fmt.Errorf("failed to add tag: %w", err)
			}
			tagIds = append(tagIds, tagId)
			continue
		} else {
			tagIds = append(tagIds, tags[0].Id)
		}
	}

	// add blog
	id, err := b.blog.Add(ctx, tx, blog)
	if err != nil {
		return nil, fmt.Errorf("failed to add blog: %w", err)
	}

	// add blogs_tags
	for _, tagId := range tagIds {
		_, err := b.blog.AddBlogTag(ctx, tx, id, tagId)
		if err != nil {
			return nil, fmt.Errorf("failed to add blogs_tags: %w", err)
		}
	}

	/// commit
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newBlog, err := b.GetBlog(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}

	return newBlog, nil
}

func (b *BlogService) ListBlog(ctx context.Context, option options.ListBlogOptions) ([]*models.Blog, error) {
	blogs, err := b.blog.List(ctx, b.db, option)
	if err != nil {
		return nil, fmt.Errorf("failed to list blog: %w", err)
	}
	return blogs, err
}

func (b *BlogService) GetBlog(ctx context.Context, id models.BlogId) (*models.Blog, error) {
	blog, err := b.blog.Get(ctx, b.db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}
	return blog, nil
}

func (b *BlogService) DeleteBlog(ctx context.Context, id models.BlogId) error {
	err := b.blog.Delete(ctx, b.db, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func (b *BlogService) PutBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// add tags
	var tagIds []models.TagId
	for _, tag := range blog.Tags {
		tags, err := b.blog.SelectTags(ctx, tx, tag)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert tag: %w", err)
		}
		if len(tags) == 0 {
			tagId, err := b.blog.AddTag(ctx, tx, tag)
			if err != nil {
				return nil, fmt.Errorf("failed to add tag: %w", err)
			}
			tagIds = append(tagIds, tagId)
			continue
		} else {
			tagIds = append(tagIds, tags[0].Id)
		}
	}

	// put blog
	id, err := b.blog.Put(ctx, tx, blog)
	if err != nil {
		return nil, fmt.Errorf("failed to put blog: %w", err)
	}

	// add blogs_tags
	for _, tagId := range tagIds {
		_, err := b.blog.AddBlogTag(ctx, tx, id, tagId)
		if err != nil {
			return nil, fmt.Errorf("failed to add blogs_tags: %w", err)
		}
	}

	/// commit
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newBlog, err := b.GetBlog(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}

	return newBlog, nil
}

func (s *BlogService) Export(ctx context.Context) error {
	return nil
}
