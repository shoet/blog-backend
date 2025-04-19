package models

type CommentId int64

type Comment struct {
	CommentId CommentId `json:"comment_id"`
	BlogId    BlogId    `json:"blog_id"`
	ClientId  *string   `json:"client_id,omitempty"`
	UserId    *UserId   `json:"user_id,omitempty"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	IsDeleted bool      `json:"is_deleted"`
	Created   int64     `json:"created"`
	Modified  int64     `json:"modified"`
}
