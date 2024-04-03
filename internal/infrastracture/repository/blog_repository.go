package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/doug-martin/goqu/v9"
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
	sql, params, err := goqu.
		Insert("blogs").
		Cols("author_id", "title", "content", "description", "thumbnail_image_file_name", "is_public").
		Vals(goqu.Vals{
			blog.AuthorId, blog.Title, blog.Content, blog.Description,
			blog.ThumbnailImageFileName, blog.IsPublic,
		}).
		Returning("id").
		ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build sql: %w", err)
	}
	row := tx.QueryRowxContext(ctx, sql, params...)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to insert blog: %w", row.Err())
	}
	var id models.BlogId
	if err := row.Scan(&id); err != nil {
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
	sql, params, err := goqu.
		Select("blog_id", goqu.C("name").As("tag")).
		From("blogs_tags").
		Join(
			goqu.T("tags"),
			goqu.On(goqu.Ex{"blogs_tags.tag_id": goqu.I("tags.id")}),
		).
		Where(goqu.Ex{"blog_id": blogId}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var tagResult []*BlogTag
	if err := tx.SelectContext(ctx, &tagResult, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to select blogs_tags: %w", err)
	}
	return tagResult, nil
}

func (r *BlogRepository) List(
	ctx context.Context, tx infrastracture.TX, option *options.ListBlogOptions,
) ([]*models.Blog, error) {
	latest := option.Limit
	builder := goqu.
		Select(
			"id", "author_id", "title", "description",
			"thumbnail_image_file_name", "is_public", "created", "modified",
		).
		From("blogs").
		Order(goqu.I("id").Desc()). // 連番なのでPKでソートする
		Limit(uint(latest))
	if option.IsPublic {
		builder = builder.Where(goqu.Ex{"is_public": true})
	}
	if option.OffsetBlogId != nil {
		builder = builder.Where(goqu.Ex{"id": goqu.Op{"gt": option.OffsetBlogId}})
	}
	sql, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	type data struct {
		models.Blog
	}
	var temp []data
	if err := tx.SelectContext(ctx, &temp, sql, params...); err != nil {
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

// ListByTagはタグ名を持つブログを検索する
func (r *BlogRepository) ListByTag(
	ctx context.Context, tx infrastracture.TX, tag string, option *options.ListBlogOptions,
) (models.Blogs, error) {
	// TODO
	// 開始のblogID以降を検索する
	// n件でlimitする
	builder := goqu.
		From("blogs_tags").
		Join(
			goqu.T("tags"),
			goqu.On(goqu.Ex{"blogs_tags.tag_id": goqu.I("tags.id")}),
		).
		Where(goqu.Ex{"tags.name": tag}).
		Select("blogs_tags.blog_id", "blogs_tags.tag_id", "tags.name").
		As("b_t")
	builder = goqu.
		From("blogs").
		Join(
			builder,
			goqu.On(goqu.Ex{"blogs.id": goqu.I("b_t.blog_id")}),
		).
		Select("blogs.*")
	if option.IsPublic {
		builder = builder.Where(goqu.Ex{"is_public": true})
	}
	sql, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var blogs models.Blogs
	if err := tx.SelectContext(ctx, &blogs, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to SelectContext: %w", err)
	}
	if len(blogs) == 0 {
		return []*models.Blog{}, nil
	}
	return blogs, nil
}

func (r *BlogRepository) ListByKeyword(
	ctx context.Context, tx infrastracture.TX, keyword string, option *options.ListBlogOptions,
) (models.Blogs, error) {
	builder := goqu.
		From("blogs").
		Where(goqu.ExOr{
			"title":       goqu.Op{"like": "%" + keyword + "%"},
			"description": goqu.Op{"like": "%" + keyword + "%"},
		}).
		Select("blogs.*")
	if option.IsPublic {
		builder = builder.Where(goqu.Ex{"is_public": true})
	}
	sql, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var blogs models.Blogs
	if err := tx.SelectContext(ctx, &blogs, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to SelectContext: %w", err)
	}
	if len(blogs) == 0 {
		return []*models.Blog{}, nil
	}
	return blogs, nil
}

func (r *BlogRepository) Get(
	ctx context.Context, tx infrastracture.TX, id models.BlogId,
) (*models.Blog, error) {
	sql, params, err := goqu.
		Select("id", "author_id", "title", "content", "description",
			"thumbnail_image_file_name", "is_public", "created", "modified",
		).
		From("blogs").
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var blogs []*models.Blog
	if err := tx.SelectContext(ctx, &blogs, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to select blog: %w", err)
	}
	if len(blogs) == 0 {
		return nil, nil
	}
	subSelectBlogsTags := goqu.Select("tag_id").From("blogs_tags").Where(goqu.Ex{"blog_id": id})
	sql, params, err = goqu.
		From(subSelectBlogsTags.As("b_t")).
		LeftOuterJoin(
			goqu.T("tags"),
			goqu.On(goqu.Ex{"b_t.tag_id": goqu.I("tags.id")}),
		).
		Select("tags.name").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}
	var tags []string
	if err := tx.SelectContext(ctx, &tags, sql, params...); err != nil {
		return nil, fmt.Errorf("failed to select tag: %w", err)
	}
	blogs[0].Tags = tags
	return blogs[0], nil
}

func (r *BlogRepository) Delete(ctx context.Context, tx infrastracture.TX, id models.BlogId) error {
	sql, params, err := goqu.
		Delete("blogs").
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build sql: %w", err)
	}
	if _, err := tx.ExecContext(ctx, sql, params...); err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func (r *BlogRepository) Put(
	ctx context.Context, tx infrastracture.TX, blog *models.Blog,
) (models.BlogId, error) {
	now := r.Clocker.Now()
	blog.Modified = uint(now.Unix())
	sql, params, err := goqu.
		Update("blogs").
		Set(goqu.Record{
			"author_id":                 blog.AuthorId,
			"title":                     blog.Title,
			"content":                   blog.Content,
			"description":               blog.Description,
			"thumbnail_image_file_name": blog.ThumbnailImageFileName,
			"is_public":                 blog.IsPublic,
			"modified":                  blog.Modified,
		}).
		Where(goqu.Ex{"id": blog.Id}).
		ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build sql: %w", err)
	}
	if _, err := tx.ExecContext(ctx, sql, params...); err != nil {
		return 0, fmt.Errorf("failed to update blog: %w", err)
	}
	return blog.Id, nil
}

func (r *BlogRepository) AddBlogTag(
	ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId,
) (int64, error) {
	sql, params, err := goqu.
		Insert("blogs_tags").
		Cols("blog_id", "tag_id").
		Vals(goqu.Vals{blogId, tagId}).
		OnConflict(
			goqu.DoUpdate(
				"blog_id, tag_id",
				goqu.Record{"blog_id": blogId, "tag_id": tagId}),
		).
		Returning("id").
		ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build sql: %w", err)
	}
	row := tx.QueryRowxContext(ctx, sql, params...)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to insert blogs_tags: %w", row.Err())
	}
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return id, nil
}

/*
そのブログと同じタグを使っている他のブログを取得する
*/
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
			blog_id = $1
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
		blog_id = $1
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
		blog_id = $1
		AND tag_id = $2
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
		name = $1
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
		($1)
	RETURNING id
	;
	`
	row := tx.QueryRowxContext(ctx, sql, tag)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to insert tags: %w", row.Err())
	}
	var tagId models.TagId
	if err := row.Scan(&tagId); err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return tagId, nil
}

func (r *BlogRepository) DeleteTag(
	ctx context.Context, tx infrastracture.TX, tagId models.TagId,
) error {
	sql := `
	DELETE FROM	
		tags
	WHERE 
		id = $1
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
	LIMIT $1
	;
	`
	var tags []*models.Tag
	if err := tx.SelectContext(ctx, &tags, sql, option.Limit); err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return tags, nil
}
