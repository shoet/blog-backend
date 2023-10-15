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
	id, err := s.blog.Add(ctx, s.db, blog)
	if err != nil {
		return fmt.Errorf("failed to add blog: %w", err)
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
