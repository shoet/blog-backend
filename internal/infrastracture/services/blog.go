package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/options"
)

func NewBlogService(db *sqlx.DB, blog BlogRepository) *BlogService {
	return &BlogService{
		db:   db,
		blog: blog,
	}
}

type BlogService struct {
	db   *sqlx.DB
	blog BlogRepository
}

func (b *BlogService) SelectTag(
	ctx context.Context, db repository.Execer, tag string,
) (*models.Tag, error) {
	tags, err := b.blog.SelectTags(ctx, db, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to select tag: %w", err)
	}
	if len(tags) == 0 {
		return nil, nil
	}
	return tags[0], nil
}

func (s *BlogService) ListTags(ctx context.Context, option options.ListTagsOptions) ([]*models.Tag, error) {
	tags, err := s.blog.ListTags(ctx, s.db, option)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	return tags, nil
}

func (s *BlogService) Export(ctx context.Context) error {
	return nil
}

func (s *BlogService) Validate(ctx context.Context, userId models.UserId, blog *models.Blog) error {
	if userId != blog.AuthorId {
		return fmt.Errorf("blog.AuthorId is invalid")
	}
	return nil
}
