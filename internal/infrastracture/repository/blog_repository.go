package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
)

type BlogRepository struct {
	Clocker clocker.Clocker
}

func NewBlogRepository(clocker clocker.Clocker) *BlogRepository {
	return &BlogRepository{
		Clocker: clocker,
	}
}

func (r *BlogRepository) Add(ctx context.Context, tx infrastracture.TX, blog *models.Blog) (models.BlogId, error) {
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
	res, err := tx.ExecContext(
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
	ctx context.Context, tx infrastracture.TX, option options.ListBlogOptions,
) ([]*models.Blog, error) {

	latest := 10
	if option.Limit != nil {
		latest = int(*option.Limit)
	}
	isPublic := ""
	if option.IsPublic {
		isPublic = "WHERE is_public = 1"
	}
	sql := `
	SELECT
		id, author_id, title, description, thumbnail_image_file_name, 
		tags.tags as tags, is_public, created, modified
	FROM (
		SELECT
			id, author_id, title, description, thumbnail_image_file_name, is_public, created, modified
		FROM
			blogs
		` + isPublic + ` 
		ORDER BY 
			created DESC
		LIMIT 
			?
	) AS blogs
	LEFT OUTER JOIN (
		SELECT
			blogs_tags.blog_id
			-- , GROUP_CONCAT(tags.name, ',') as tags -- for sqlite3
			, GROUP_CONCAT(tags.name) as tags -- for mysql
		FROM blogs_tags
		JOIN tags
			ON blogs_tags.tag_id = tags.id
		GROUP BY blogs_tags.blog_id
		-- TODO: 将来的に遅くなる チューニング
	) AS tags
		ON blogs.id = tags.blog_id
	;
	`
	type data struct {
		Id                     models.BlogId `json:"id" db:"id"`
		Title                  string        `json:"title" db:"title"`
		Description            string        `json:"description" db:"description"`
		Content                string        `json:"content,omitempty" db:"content"`
		AuthorId               models.UserId `json:"authorId" db:"author_id"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName" db:"thumbnail_image_file_name"`
		IsPublic               bool          `json:"isPublic" db:"is_public"`
		Tags                   *string       `json:"tags" db:"tags"`
		Created                time.Time     `json:"created" db:"created"`
		Modified               time.Time     `json:"modified" db:"modified"`
	}
	var temp []data
	if err := tx.SelectContext(ctx, &temp, sql, latest); err != nil {
		return nil, fmt.Errorf("failed to select blogs: %w", err)
	}
	var blogs []*models.Blog
	for _, t := range temp {
		var tags []string
		if t.Tags != nil {
			tags = strings.Split(*t.Tags, ",")
		}
		blogs = append(blogs, &models.Blog{
			Id:                     t.Id,
			Title:                  t.Title,
			Description:            t.Description,
			Content:                t.Content,
			AuthorId:               t.AuthorId,
			ThumbnailImageFileName: t.ThumbnailImageFileName,
			IsPublic:               t.IsPublic,
			Tags:                   tags,
			Created:                t.Created,
			Modified:               t.Modified,
		})
	}
	return blogs, nil
}

func (r *BlogRepository) Get(
	ctx context.Context, tx infrastracture.TX, id models.BlogId,
) (*models.Blog, error) {
	sqlBlog := `
	SELECT
		id, author_id, title, content, description,
		thumbnail_image_file_name, is_public, created, modified
	FROM
		blogs WHERE id = ?
	;
	`
	var blogs []*models.Blog
	if err := tx.SelectContext(ctx, &blogs, sqlBlog, id); err != nil {
		return nil, fmt.Errorf("failed to select blog: %w", err)
	}
	if len(blogs) == 0 {
		return nil, nil
	}

	sqlTags := `
	SELECT
		tags.name
	FROM (
		SELECT
			tag_id
		FROM blogs_tags
		WHERE blog_id = ?
	) as b_t
	LEFT OUTER JOIN tags
		ON b_t.tag_id = tags.id
	`
	var tags []string
	if err := tx.SelectContext(ctx, &tags, sqlTags, id); err != nil {
		return nil, fmt.Errorf("failed to select tag: %w", err)
	}
	blogs[0].Tags = tags
	return blogs[0], nil
}

func (r *BlogRepository) Delete(ctx context.Context, tx infrastracture.TX, id models.BlogId) error {
	sql := `
	DELETE FROM
		blogs
	WHERE 
		id = ?
	;
	`
	_, err := tx.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func (r *BlogRepository) Put(
	ctx context.Context, db Execer, blog *models.Blog,
) (models.BlogId, error) {
	sql := `
	UPDATE blogs
	SET
		author_id = ?
		, title = ?
		, content = ?
		, description = ?
		, thumbnail_image_file_name = ?
		, is_public = ?
		, modified = ?
	WHERE
		id = ?
	;
	`
	now := r.Clocker.Now()
	blog.Modified = now
	_, err := db.ExecContext(
		ctx,
		sql,
		blog.AuthorId, blog.Title, blog.Content, blog.Description,
		blog.ThumbnailImageFileName, blog.IsPublic, blog.Modified, blog.Id)
	if err != nil {
		return 0, fmt.Errorf("failed to update blog: %w", err)
	}
	return blog.Id, nil
}

func (r *BlogRepository) AddBlogTag(
	ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId,
) (int64, error) {
	sql := `
	REPLACE INTO blogs_tags
		(blog_id, tag_id)
	VALUES
		(?, ?)
	;
	`
	res, err := tx.ExecContext(ctx, sql, blogId, tagId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert blogs_tags: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return id, nil
}

func (r *BlogRepository) SelectBlogsTagsByOtherUsingBlog(
	ctx context.Context, tx infrastracture.TX, blogId models.BlogId,
) ([]*models.BlogsTags, error) {
	var result []*models.BlogsTags
	sql := `
	SELECT 
		a.blog_id
		, a.tag_id
		, tags.name
	FROM 
		blogs_tags as a
	JOIN (
		SELECT
			blog_id
			, tag_id
		FROM
			blogs_tags
		WHERE
			blog_id = ?
	) as b
	ON 
		a.tag_id = b.tag_id
		AND
			a.blog_id <> b.blog_id
	LEFT OUTER JOIN tags
		ON a.tag_id = tags.id
	`
	if err := tx.SelectContext(ctx, &result, sql, blogId); err != nil {
		return nil, fmt.Errorf("failed to select using tags: %w", err)
	}
	return result, nil
}

func (r *BlogRepository) SelectBlogsTags(
	ctx context.Context, tx infrastracture.TX, blogId models.BlogId,
) ([]*models.BlogsTags, error) {
	var result []*models.BlogsTags
	sql := `
	SELECT
		blogs_tags.blog_id
		, blogs_tags.tag_id
		, tags.name
	FROM
		blogs_tags
	LEFT OUTER JOIN tags
		ON blogs_tags.tag_id = tags.id
	WHERE
		blog_id = ?
	;
	`
	if err := tx.SelectContext(ctx, &result, sql, blogId); err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return result, nil
}

func (r *BlogRepository) DeleteBlogsTags(
	ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId,
) error {
	sql := `
	DELETE FROM
		blogs_tags
	WHERE
		blog_id = ?
		AND tag_id = ?
	;
	`
	if _, err := tx.ExecContext(ctx, sql, blogId, tagId); err != nil {
		return fmt.Errorf("failed to delete blogs_tags: %w", err)
	}
	return nil
}

func (r *BlogRepository) SelectTags(
	ctx context.Context, tx infrastracture.TX, tag string,
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
	if err := tx.SelectContext(ctx, &tags, sql, tag); err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return tags, nil
}

func (r *BlogRepository) AddTag(ctx context.Context, tx infrastracture.TX, tag string) (models.TagId, error) {
	sql := `
	INSERT INTO tags
		(name)
	VALUES
		(?)
	;
	`
	res, err := tx.ExecContext(ctx, sql, tag)
	if err != nil {
		return 0, fmt.Errorf("failed to insert tags: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return models.TagId(id), nil
}

func (r *BlogRepository) DeleteTag(
	ctx context.Context, tx infrastracture.TX, tagId models.TagId,
) error {
	sql := `
	DELETE FROM	
		tags
	WHERE 
		id = ?
	;
	`
	if _, err := tx.ExecContext(ctx, sql, tagId); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (r *BlogRepository) ListTags(
	ctx context.Context, db Queryer, option options.ListTagsOptions,
) ([]*models.Tag, error) {
	sql := `
	SELECT
		id, name
	FROM
		tags
	ORDER BY
		name
	LIMIT ?
	;
	`
	var tags []*models.Tag
	if err := db.SelectContext(ctx, &tags, sql, option.Limit); err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return tags, nil
}
