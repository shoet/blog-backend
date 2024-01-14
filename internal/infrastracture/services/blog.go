package services

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture/models"
)

type BlogService struct{}

func NewBlogService() *BlogService {
	return &BlogService{}
}

func (s *BlogService) Validate(ctx context.Context, userId models.UserId, blog *models.Blog) error {
	if userId != blog.AuthorId {
		return fmt.Errorf("blog.AuthorId is invalid")
	}
	return nil
}
