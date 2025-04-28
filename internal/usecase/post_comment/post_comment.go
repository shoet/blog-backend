package post_comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type CommentRepository interface {
	Get(ctx context.Context, tx infrastructure.TX, commentId models.CommentId) (*models.Comment, error)
	CreateComment(
		ctx context.Context,
		tx infrastructure.TX,
		blogId models.BlogId,
		userId *models.UserId,
		clientId *string,
		threadId *string,
		content string,
	) (models.CommentId, error)

	UpdateThreadId(ctx context.Context, tx infrastructure.TX, commentId models.CommentId, threadId string) error
}

type Usecase struct {
	DB                infrastructure.DB
	CommentRepository CommentRepository
}

func NewUsecase(
	db infrastructure.DB,
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
	threadCommentId *int64,
	content string,
) (models.CommentId, error) {
	if userId == nil && clientId == nil {
		return 0, fmt.Errorf("UserID or ClientID is required")
	}

	transactionProvider := infrastructure.NewTransactionProvider(u.DB)
	result, err := transactionProvider.DoInTx(ctx, func(tx infrastructure.TX) (any, error) {
		var threadId *string
		// コメントをスレッドにする場合、threadCommentIdとしてスレッド化対象のコメントIDが指定される
		// スレッドIDを発行し、関連するコメントを一つのスレッドIDに紐づける
		if threadCommentId != nil {
			commentId := models.CommentId(*threadCommentId)
			comment, err := u.CommentRepository.Get(ctx, tx, commentId)
			if err != nil {
				return 0, fmt.Errorf("failed to get comment: %w", err)
			}
			if comment.ThreadId != nil {
				threadId = comment.ThreadId
			} else {
				tid := uuid.New().String()
				threadId = &tid
			}
			if err := u.CommentRepository.UpdateThreadId(ctx, tx, models.CommentId(*threadCommentId), *threadId); err != nil {
				return 0, fmt.Errorf("failed to update thread id: %w", err)
			}
		}
		commentId, err := u.CommentRepository.CreateComment(ctx, tx, blogId, userId, clientId, threadId, content)
		if err != nil {
			return 0, fmt.Errorf("failed to create comment: %w", err)
		}
		return commentId, nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}
	commentId, ok := result.(models.CommentId)
	if !ok {
		return 0, fmt.Errorf("failed to cast result to CommentId")
	}
	return commentId, nil
}
