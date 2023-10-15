package store

import (
	"context"
	"fmt"

	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/interfaces"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogRepository struct {
	Clocker clocker.Clocker
}

func (r *BlogRepository) Add(ctx context.Context, db interfaces.Execer, blog *models.Blog) (models.BlogId, error) {
	sql := `
	INSERT INTO blogs
		(author_id, title, content, description, thumbnail_image_file_name, is_public, created, modified)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	;
	`
	now := r.Clocker.Now()
	blog.Created = now
	blog.Modified = now
	res, err := db.ExecContext(
		ctx,
		sql,
		blog.AuthorId, blog.Title, blog.Content, blog.Description,
		blog.ThumbnailImageFileName, blog.IsPublic, blog.Created, blog.Modified)
	if err != nil {
		return 0, fmt.Errorf("failed to insert blog: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return models.BlogId(id), nil
}

func (r *BlogRepository) List(
	ctx context.Context, db interfaces.Queryer, option options.ListBlogOptions,
) ([]*models.Blog, error) {
	sql := `
	SELECT
		id, author_id, title, content, description,thumbnail_image_file_name, is_public, created, modified
	FROM
		blogs
	;
	`
	var blogs []*models.Blog
	if err := db.SelectContext(ctx, &blogs, sql); err != nil {
		return nil, fmt.Errorf("failed to select blogs: %w", err)
	}
	return blogs, nil
}

func (r *BlogRepository) Get(
	ctx context.Context, db interfaces.Queryer, id models.BlogId,
) (*models.Blog, error) {
	sql := `
	SELECT
		id, author_id, title, content, description,
		thumbnail_image_file_name, is_public, created, modified
	FROM
		blogs
	where
		id = ?
	;
	`
	var blog []*models.Blog
	if err := db.SelectContext(ctx, &blog, sql, id); err != nil {
		return nil, fmt.Errorf("failed to select blog: %w", err)
	}
	if len(blog) == 0 {
		return nil, nil
	}
	return blog[0], nil
}

func (r *BlogRepository) Delete(ctx context.Context, db interfaces.Execer, id models.BlogId) error {
	return nil
}
func (r *BlogRepository) Put(ctx context.Context, db interfaces.Execer, blog *models.Blog) error {
	return nil
}

func (r *BlogRepository) AddBlogTag(
	ctx context.Context, db interfaces.Execer, blogId models.BlogId, tagId models.TagId,
) (int64, error) {
	sql := `
	REPLACE INTO blogs_tags
		(blog_id, tag_id)
	VALUES
		(?, ?)
	;
	`
	res, err := db.ExecContext(ctx, sql, blogId, tagId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert blogs_tags: %w", err)
	}
	id, err := res.LastInsertId()
	return id, nil
}

func (r *BlogRepository) SelectTags(
	ctx context.Context, db interfaces.Queryer, tag string,
) ([]*models.Tag, error) {
	sql := `
	SELECT
		id, name
	FROM
		tags
	WHERE
		name = ?
	;
	`
	var tags []*models.Tag
	if err := db.SelectContext(ctx, &tags, sql, tag); err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return tags, nil
}

func (r *BlogRepository) AddTag(ctx context.Context, db interfaces.Execer, tag string) (models.TagId, error) {
	sql := `
	INSERT INTO tags
		(name)
	VALUES
		(?)
	;
	`
	res, err := db.ExecContext(ctx, sql, tag)
	if err != nil {
		return 0, fmt.Errorf("failed to insert tags: %w", err)
	}
	id, err := res.LastInsertId()
	return models.TagId(id), nil

}
