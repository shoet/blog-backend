package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"

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
	blog.Created = uint(now.Unix())
	blog.Modified = uint(now.Unix())
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

type BlogTag struct {
	BlogId models.BlogId `db:"blog_id"`
	Tag    string        `db:"tag"`
}

// Blogに紐づくTagを取得する
func (r *BlogRepository) WithBlogTags(
	ctx context.Context, tx infrastracture.TX, blogId models.BlogId,
) ([]*BlogTag, error) {
	sql := `
	SELECT
		blogs_tags.blog_id, tags.name as tag
	FROM 
		blogs_tags
	JOIN
		tags
	ON
		blogs_tags.tag_id = tags.id
	WHERE
		blog_id = ?
	;	
	`
	var tagResult []*BlogTag
	if err := tx.SelectContext(ctx, &tagResult, sql, blogId); err != nil {
		return nil, fmt.Errorf("failed to select blogs_tags: %w", err)
	}
	return tagResult, nil
}

func (r *BlogRepository) List(
	ctx context.Context, tx infrastracture.TX, option *options.ListBlogOptions,
) ([]*models.Blog, error) {
	latest := option.Limit
	isPublic := ""
	if option.IsPublic {
		isPublic = "WHERE is_public = true"
	}
	sql := `
	SELECT
		id, author_id, title, description, 
		thumbnail_image_file_name, is_public, created, modified
	FROM
		blogs
	` + isPublic + ` 
	ORDER BY 
		id DESC -- 連番なのでPKでソートする
	LIMIT 
		$1
	;
	`
	type data struct {
		models.Blog
	}
	var temp []data
	if err := tx.SelectContext(ctx, &temp, sql, latest); err != nil {
		return nil, fmt.Errorf("failed to select blogs: %w", err)
	}
	var blogs []*models.Blog
	for _, t := range temp {
		blogTag, err := r.WithBlogTags(ctx, tx, t.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to select blogs_tags: %w", err)
		}
		tags := make([]string, 0, len(blogTag))
		for _, t := range blogTag {
			tags = append(tags, t.Tag)
		}
		// タグを昇順にソート
		sort.SliceStable(tags, func(i, j int) bool {
			return strings.Compare(tags[i], tags[j]) < 0
		})
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
	ctx context.Context, tx infrastracture.TX, blog *models.Blog,
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
	blog.Modified = uint(now.Unix())
	_, err := tx.ExecContext(
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
	ctx context.Context, tx infrastracture.TX, option options.ListTagsOptions,
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
	if err := tx.SelectContext(ctx, &tags, sql, option.Limit); err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return tags, nil
}
