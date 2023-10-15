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

func (s *BlogService) AddBlog(ctx context.Context, blog *models.Blog) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// add tags
	var tagIds []models.TagId
	for _, tag := range blog.Tags {
		tags, err := s.blog.SelectTags(ctx, tx, tag)
		if err != nil {
			return fmt.Errorf("failed to upsert tag: %w", err)
		}
		if len(tags) == 0 {
			tagId, err := s.blog.AddTag(ctx, tx, tag)
			if err != nil {
				return fmt.Errorf("failed to add tag: %w", err)
			}
			tagIds = append(tagIds, tagId)
			continue
		} else {
			tagIds = append(tagIds, tags[0].Id)
		}
	}

	// add blog
	id, err := s.blog.Add(ctx, tx, blog)
	if err != nil {
		return fmt.Errorf("failed to add blog: %w", err)
	}

	// add blogs_tags
	for _, tagId := range tagIds {
		_, err := s.blog.AddBlogTag(ctx, tx, id, tagId)
		if err != nil {
			return fmt.Errorf("failed to add blogs_tags: %w", err)
		}
	}

	/// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	blog.Id = id
	return nil
}

func (s *BlogService) ListBlog(ctx context.Context, option options.ListBlogOptions) ([]*models.Blog, error) {
	blogs, err := s.blog.List(ctx, s.db, option)
	if err != nil {
		return nil, fmt.Errorf("failed to list blog: %w", err)
	}
	return blogs, err
}

func (s *BlogService) GetBlog(ctx context.Context, id models.BlogId) (*models.Blog, error) {
	blog, err := s.blog.Get(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}
	return blog, nil
}

func (s *BlogService) DeleteBlog(ctx context.Context, id models.BlogId) error {
	err := s.blog.Delete(ctx, s.db, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func (s *BlogService) PutBlog(ctx context.Context, blog *models.Blog) error {
	err := s.blog.Put(ctx, s.db, blog)
	if err != nil {
		return fmt.Errorf("failed to put blog: %w", err)
	}
	return nil
}
