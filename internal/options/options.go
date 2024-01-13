package options

import "github.com/shoet/blog/internal/infrastracture/models"

// TODO: リファクタリング
type ListBlogOptions struct {
	// AuthorId *models.UserId // 未使用
	Tags     []models.TagId
	IsPublic bool
	Limit    *int64 // default 10
}

type ListTagsOptions struct {
	Limit int
}
