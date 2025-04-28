package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
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
	tx infrastructure.TX,
	blogId models.BlogId,
	userId *models.UserId,
	clientId *string,
	threadId *string,
	content string,
) (models.CommentId, error) {
	builder := goqu.
		Insert("comments").
		Cols("blog_id", "client_id", "user_id", "thread_id", "content", "created", "modified").
		Returning("comment_id").
		Rows(
			goqu.Record{
				"blog_id": blogId, "client_id": clientId, "thread_id": threadId, "user_id": userId, "content": content,
				"created": r.Clocker.Now(), "modified": r.Clocker.Now(),
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

func (r *CommentRepository) Get(
	ctx context.Context,
	tx infrastructure.TX,
	commentId models.CommentId,
) (*models.Comment, error) {
	builder := goqu.
		Select("comment_id", "blog_id", "client_id", "user_id", "thread_id", "content", "is_edited", "is_deleted", "created", "modified").
		From("comments").
		Where(goqu.Ex{"comment_id": commentId})
	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	row := tx.QueryRowxContext(ctx, query, params...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to get comment: %w", row.Err())
	}
	comment := &models.Comment{}
	if err := row.StructScan(comment); err != nil {
		return nil, fmt.Errorf("failed to scan comment: %w", err)
	}
	return comment, nil
}

func (r *CommentRepository) UpdateThreadId(
	ctx context.Context,
	tx infrastructure.TX,
	commentId models.CommentId,
	threadId string,
) error {
	builder := goqu.
		Update("comments").
		Set(goqu.Record{"thread_id": threadId}).
		Where(goqu.Ex{"comment_id": commentId})
	query, params, err := builder.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	if _, err := tx.ExecContext(ctx, query, params...); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (r *CommentRepository) GetByBlogId(
	ctx context.Context,
	tx infrastructure.TX,
	blogId models.BlogId,
	excludeDeleted bool,
) ([]*models.Comment, error) {
	builder := goqu.
		Select("comment_id", "blog_id", "client_id", "user_id", "thread_id", "content", "is_edited", "is_deleted", "created", "modified").
		From("comments").
		Where(goqu.Ex{"blog_id": blogId}).
		Where(goqu.Ex{"is_deleted": !excludeDeleted}).
		Order(goqu.I("created").Asc())
	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	rows, err := tx.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", rows.Err())
	}
	comments := make([]*models.Comment, 0, 0)
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.StructScan(&comment); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
