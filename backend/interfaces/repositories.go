package interfaces

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogRepository interface {
	Add(ctx context.Context, db *sqlx.DB, blog *models.Blog) (models.BlogId, error)
	List(ctx context.Context, db *sqlx.DB, option options.ListBlogOptions) ([]*models.Blog, error)
	Delete(ctx context.Context, db *sqlx.DB, id models.BlogId) error
	Put(ctx context.Context, db *sqlx.DB, blog *models.Blog) error
}

type UserRepository interface {
	Add(ctx context.Context, db *sqlx.DB, user *models.User) (models.UserId, error)
	Get(ctx context.Context, db *sqlx.DB, id models.UserId) (*models.User, error)
	GetByEmail(ctx context.Context, db *sqlx.DB, email string) (*models.User, error)
	Delete(ctx context.Context, db *sqlx.DB, id models.UserId) error
	Put(ctx context.Context, db *sqlx.DB, user *models.User) error
}
