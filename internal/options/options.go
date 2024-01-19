package options

type ListBlogOptions struct {
	IsPublic bool
	Limit    int64
}

func NewListBlogOptions(isPublic *bool, limit *int64) *ListBlogOptions {
	option := new(ListBlogOptions)
	if isPublic == nil {
		option.IsPublic = false // Default
	} else {
		option.IsPublic = *isPublic
	}
	if limit == nil {
		option.Limit = 10 // Default
	} else {
		option.Limit = *limit
	}
	return option
}

type ListTagsOptions struct {
	Limit int
}
