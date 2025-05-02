package models

import "time"

type CommentId int64

type Comment struct {
	CommentId CommentId `json:"commentId" db:"comment_id"`
	BlogId    BlogId    `json:"blogId" db:"blog_id"`
	ClientId  *string   `json:"clientId,omitempty" db:"client_id"`
	UserId    *UserId   `json:"userId,omitempty" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	IsEdited  bool      `json:"isEdited" db:"is_edited"`
	IsDeleted bool      `json:"isDeleted" db:"is_deleted"`
	ThreadId  *string   `json:"threadId,omitempty" db:"thread_id"`
	Created   time.Time `json:"created" db:"created"`
	Modified  time.Time `json:"modified" db:"modified"`

	Nickname           *string `json:"nickname,omitempty"`
	AvatarImageFileURL *string `json:"avatarImageFileUrl,omitempty"`
}
