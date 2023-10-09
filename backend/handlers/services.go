package handlers

import "github.com/shoet/blog/models"

type ListBlogOptions struct {
	AuthorId models.UserId
	Tags     []models.TagId
	IsPublic bool
}

type BlogService interface {
	ListBlog(option *ListBlogOptions) ([]*models.Blog, error)
	AddBlog(*models.Blog) error
	DeleteBlog(models.BlogId) error
	PutBlog(*models.Blog) error
}
