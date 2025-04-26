package models

import "slices"

type TagId int64

type Tag struct {
	Id   TagId  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Tags []*Tag

type BlogsTags struct {
	BlogId BlogId `json:"blogId" db:"blog_id"`
	TagId  TagId  `json:"tagId" db:"tag_id"`
	Name   string `json:"name" db:"name"`
}

// BlogsTagsArray は、ブログとタグのリレーションを保持する構造体のスライス
type BlogsTagsArray []*BlogsTags

// TagIds は、リレーションのスライスからタグIDのスライスを取得する関数
func (arr BlogsTagsArray) TagIds() []TagId {
	var tags []TagId
	for _, bt := range arr {
		tags = append(tags, bt.TagId)
	}
	return tags
}

// TagNames は、リレーションのスライスからタグ名のスライスを取得する関数
func (arr BlogsTagsArray) TagNames() []string {
	var tags []string
	for _, t := range arr {
		tags = append(tags, t.Name)
	}
	return tags
}

// Contains は、リレーションのスライスに指定したタグが含まれているかを判定する関数
func (arr BlogsTagsArray) Contains(tag string) bool {
	return slices.Contains(arr.TagNames(), tag)
}
