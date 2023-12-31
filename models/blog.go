package models

import "time"

type BlogId int64

type Blog struct {
	Id                     BlogId    `json:"id" db:"id"`
	Title                  string    `json:"title" db:"title"`
	Description            string    `json:"description" db:"description"`
	Content                string    `json:"content,omitempty" db:"content"`
	AuthorId               UserId    `json:"authorId" db:"author_id"`
	ThumbnailImageFileName string    `json:"thumbnailImageFileName" db:"thumbnail_image_file_name"`
	IsPublic               bool      `json:"isPublic" db:"is_public"`
	Tags                   []string  `json:"tags,omitempty" db:"tags"`
	Created                time.Time `json:"created" db:"created"`
	Modified               time.Time `json:"modified" db:"modified"`
}

type Blogs []*Blog

type UserId int64

type User struct {
	Id       UserId    `json:"id,omitempty" db:"id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email,omitempty" db:"email"`
	Password string    `json:"password,omitempty" db:"password"`
	Created  time.Time `json:"created,omitempty" db:"created"`
	Modified time.Time `json:"modified,omitempty" db:"modified"`
}

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

type BlogsTagsArray []*BlogsTags

func (arr BlogsTagsArray) TagIds() []TagId {
	var tags []TagId
	for _, bt := range arr {
		tags = append(tags, bt.TagId)
	}
	return tags
}

func (arr BlogsTagsArray) TagNames() []string {
	var tags []string
	for _, t := range arr {
		tags = append(tags, t.Name)
	}
	return tags
}
