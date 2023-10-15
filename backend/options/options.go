package options

import "github.com/shoet/blog/models"

type ListBlogOptions struct {
	AuthorId models.UserId
	Tags     []models.TagId
	IsPublic bool
}
