package models

type CommentId int64

type Comment struct {
	CommentId CommentId `json:"comment_id" db:"comment_id"`
	BlogId    BlogId    `json:"blog_id" db:"blog_id"`
	ClientId  *string   `json:"client_id,omitempty" db:"client_id"`
	UserId    *UserId   `json:"user_id,omitempty" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	IsEdited  bool      `json:"is_edited" db:"is_edited"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
	ThreadId  *string   `json:"thread_id,omitempty" db:"thread_id"`
	Created   int64     `json:"created" db:"created"`
	Modified  int64     `json:"modified" db:"modified"`
}
