package options

import "github.com/shoet/blog/models"

type ListBlogOptions struct {
	AuthorId *models.UserId
	Tags     []models.TagId
	IsPublic bool
	Limit    *int64
}

type ListTagsOptions struct {
	Limit int
}
