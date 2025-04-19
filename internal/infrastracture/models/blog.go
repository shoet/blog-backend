package models

import (
	"strings"

	"golang.org/x/exp/slices"
)

type BlogId int64

type Blog struct {
	Id                     BlogId     `json:"id" db:"id"`
	Title                  string     `json:"title" db:"title"`
	Description            string     `json:"description" db:"description"`
	Content                string     `json:"content,omitempty" db:"content"`
	AuthorId               UserId     `json:"authorId" db:"author_id"`
	ThumbnailImageFileName string     `json:"thumbnailImageFileName" db:"thumbnail_image_file_name"`
	IsPublic               bool       `json:"isPublic" db:"is_public"`
	Tags                   []string   `json:"tags,omitempty" db:"tags"`
	Comments               []*Comment `json:"comments" db:"comments"`
	Created                uint       `json:"created" db:"created"`
	Modified               uint       `json:"modified" db:"modified"`
}

func (blog *Blog) HavingTag(tag string) bool {
	if slices.Contains(blog.Tags, tag) {
		return true
	}
	return false
}

func (blog *Blog) HavingKeyword(keyword string) bool {
	// タイトルか概要にキーワードが含まれている
	keywordLower := strings.ToLower(keyword)
	titleLower := strings.ToLower(blog.Title)
	descriptionLower := strings.ToLower(blog.Description)
	return strings.Contains(titleLower, keywordLower) || strings.Contains(descriptionLower, keywordLower)
}

type Blogs []*Blog

func (blogs Blogs) FilterByTag(tag string) Blogs {
	var result Blogs
	for _, blog := range blogs {
		if blog.HavingTag(tag) {
			result = append(result, blog)
		}
	}
	return result
}

func (blogs Blogs) FilterByKeyword(keyword string) Blogs {
	var result Blogs
	for _, blog := range blogs {
		if blog.HavingKeyword(keyword) {
			result = append(result, blog)
		}
	}
	return result
}

func (blogs Blogs) ToSlice() []*Blog {
	result := make([]*Blog, 0, len(blogs))
	for _, blog := range blogs {
		result = append(result, blog)
	}
	return result
}
