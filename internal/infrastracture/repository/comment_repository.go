package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type CommentRepository struct {
	Clocker clocker.Clocker
}

func NewCommentRepository(clocker clocker.Clocker) *CommentRepository {
	return &CommentRepository{
		Clocker: clocker,
	}
}

func (r *CommentRepository) CreateComment(
	ctx context.Context,
	tx infrastracture.TX,
	blogId models.BlogId,
	userId *models.UserId,
	clientId *string,
	content string,
) (models.CommentId, error) {
	builder := goqu.
		Insert("comments").
		Cols("blog_id", "client_id", "user_id", "content", "created", "modified").
		Returning("comment_id").
		Rows(
			goqu.Record{
				"blog_id": blogId, "client_id": clientId, "user_id": userId, "content": content,
				"created": r.Clocker.Now().Unix(), "modified": r.Clocker.Now().Unix(),
			},
		)
	query, params, err := builder.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}
	row := tx.QueryRowxContext(ctx, query, params...)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to insert comment: %w", row.Err())
	}
	var commentId models.CommentId
	if err := row.Scan(&commentId); err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return commentId, nil
}
