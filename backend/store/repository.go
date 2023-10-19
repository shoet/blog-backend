package store

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	if option.AuthorId == nil {
		return nil, fmt.Errorf("author id is required")
	}
	latest := 10
	if option.Limit != nil {
		latest = int(*option.Limit)
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
		WHERE
			author_id = ?
		ORDER BY 
			created DESC
		LIMIT 
			?
	) AS blogs
	LEFT OUTER JOIN (
		SELECT
			blogs_tags.blog_id
			, GROUP_CONCAT(tags.name, ',') as tags
		FROM blogs_tags
		JOIN tags
			ON blogs_tags.tag_id = tags.id
		GROUP BY blogs_tags.blog_id
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
	if err := db.SelectContext(ctx, &temp, sql, option.AuthorId, latest); err != nil {
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
	ctx context.Context, db interfaces.Queryer, id models.BlogId,
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
	if err := db.SelectContext(ctx, &blogs, sqlBlog, id); err != nil {
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
	if err := db.SelectContext(ctx, &tags, sqlTags, id); err != nil {
		return nil, fmt.Errorf("failed to select tag: %w", err)
	}
	blogs[0].Tags = tags
	return blogs[0], nil
}

func (r *BlogRepository) Delete(ctx context.Context, db interfaces.Execer, id models.BlogId) error {
	sql := `
	DELETE FROM
		blogs
	WHERE 
		id = ?
	;
	`
	_, err := db.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}
func (r *BlogRepository) Put(
	ctx context.Context, db interfaces.Execer, blog *models.Blog,
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

type UserRepository struct {
	Clocker clocker.Clocker
}

func (u *UserRepository) GetByEmail(
	ctx context.Context, db interfaces.Queryer, email string,
) (*models.User, error) {
	sql := `
	SELECT
		id, name, email, password, created, modified
	FROM users
	WHERE email = ?
	;
	`
	fmt.Println(email)
	var users []*models.User
	if err := db.SelectContext(ctx, &users, sql, email); err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return users[0], nil
}
