package interfaces

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type Queryer interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

var _ Execer = (*sqlx.Tx)(nil)
var _ Execer = (*sqlx.DB)(nil)
var _ Queryer = (*sqlx.Tx)(nil)
var _ Queryer = (*sqlx.DB)(nil)

type BlogRepository interface {
	Add(ctx context.Context, db Execer, blog *models.Blog) (models.BlogId, error)
	List(ctx context.Context, db Queryer, option options.ListBlogOptions) ([]*models.Blog, error)
	Get(ctx context.Context, db Queryer, id models.BlogId) (*models.Blog, error)
	Delete(ctx context.Context, db Execer, id models.BlogId) error
	Put(ctx context.Context, db Execer, blog *models.Blog) error
	AddBlogTag(ctx context.Context, db Execer, blogId models.BlogId, tagId models.TagId) (int64, error)
	SelectTags(ctx context.Context, db Queryer, tag string) ([]*models.Tag, error)
	AddTag(ctx context.Context, db Execer, tag string) (models.TagId, error)
}

type UserRepository interface {
	Add(ctx context.Context, db Execer, user *models.User) (models.UserId, error)
	Get(ctx context.Context, db Queryer, id models.UserId) (*models.User, error)
	GetByEmail(ctx context.Context, db Queryer, email string) (*models.User, error)
	Delete(ctx context.Context, db Execer, id models.UserId) error
	Put(ctx context.Context, db Execer, user *models.User) error
}
