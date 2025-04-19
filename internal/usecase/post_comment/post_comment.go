package post_comment

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type CommentRepository interface {
	CreateComment(
		ctx context.Context,
		tx infrastracture.TX,
		blogId models.BlogId,
		userId *models.UserId,
		clientId *string,
		content string,
	) (models.CommentId, error)
}

type Usecase struct {
	DB                infrastracture.DB
	CommentRepository CommentRepository
}

func NewUsecase(
	db infrastracture.DB,
	commentRepository CommentRepository,
) *Usecase {
	return &Usecase{
		DB:                db,
		CommentRepository: commentRepository,
	}
}

func (u *Usecase) Run(
	ctx context.Context,
	blogId models.BlogId,
	userId *models.UserId,
	clientId *string,
	content string,
) (models.CommentId, error) {
	if userId == nil && clientId == nil {
		return 0, fmt.Errorf("UserID or ClientID is required")
	}
	commentId, err := u.CommentRepository.CreateComment(ctx, u.DB, blogId, userId, clientId, content)
	if err != nil {
		return 0, fmt.Errorf("failed to create comment")
	}
	return commentId, nil
}
