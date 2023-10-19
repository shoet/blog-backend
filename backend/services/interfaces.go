package services

import (
	"context"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
	"github.com/shoet/blog/store"
)

//go:generate go run github.com/matryer/moq -out service_moq.go . JWTer KVSer

type BlogRepository interface {
	Add(ctx context.Context, db store.Execer, blog *models.Blog) (models.BlogId, error)
	List(ctx context.Context, db store.Queryer, option options.ListBlogOptions) ([]*models.Blog, error)
	Get(ctx context.Context, db store.Queryer, id models.BlogId) (*models.Blog, error)
	Delete(ctx context.Context, db store.Execer, id models.BlogId) error
	Put(ctx context.Context, db store.Execer, blog *models.Blog) (models.BlogId, error)
	AddBlogTag(ctx context.Context, db store.Execer, blogId models.BlogId, tagId models.TagId) (int64, error)
	SelectTags(ctx context.Context, db store.Queryer, tag string) ([]*models.Tag, error)
	AddTag(ctx context.Context, db store.Execer, tag string) (models.TagId, error)
}

type UserRepository interface {
	// Add(ctx context.Context, db store.Execer, user *models.User) (models.UserId, error)
	// Get(ctx context.Context, db store.Queryer, id models.UserId) (*models.User, error)
	GetByEmail(ctx context.Context, db store.Queryer, email string) (*models.User, error)
	// Delete(ctx context.Context, db store.Execer, id models.UserId) error
	// Put(ctx context.Context, db store.Execer, user *models.User) error
}

type JWTer interface {
	GenerateToken(ctx context.Context, u *models.User) (string, error)
	VerifyToken(ctx context.Context, token string) (models.UserId, error)
}

type KVSer interface {
	Save(ctx context.Context, key string, value string) error
	Load(ctx context.Context, key string) (string, error)
}
