package options

type ListBlogOptions struct {
	Tag      *string
	KeyWord  *string
	IsPublic bool
	Limit    int64
}

func NewListBlogOptions(
	tag *string, keyWord *string, isPublic *bool, limit *int64,
) *ListBlogOptions {
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
	option.Tag = tag
	option.KeyWord = keyWord
	return option
}

type ListTagsOptions struct {
	Limit int
}
