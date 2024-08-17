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

// BlogRepositoryOffset はブログリポジトリをオフセットベースでのページネーションを実装した構造体です。
// BlogRepository を埋め込んでいます。
type BlogRepositoryOffset struct {
	*BlogRepository
}

func NewBlogRepositoryOffset(clocker clocker.Clocker) *BlogRepositoryOffset {
	return &BlogRepositoryOffset{
		BlogRepository: NewBlogRepository(clocker),
	}
}

// buildOffset は指定されたpage, limitに対するオフセットを生成します。
func (r *BlogRepositoryOffset) buildOffset(page int64, limit int64) int64 {
	// offset := page * limit
	offset := (page - 1) * limit
	return offset
}

func (r *BlogRepositoryOffset) List(
	ctx context.Context, tx infrastracture.TX, option *options.ListBlogOptions,
) (models.Blogs, error) {
	builder := goqu.
		Select(
			"id", "author_id", "title", "description",
			"thumbnail_image_file_name", "is_public", "created", "modified",
		).
		From("blogs").
		Order(goqu.I("id").Desc()).
		Limit(uint(option.Limit))
	if option.IsPublic {
		builder = builder.Where(goqu.Ex{"is_public": true})
	}
	if option.Page != nil {
		offset := r.buildOffset(*option.Page, option.Limit)
		builder = builder.Offset(uint(offset))
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
	blogs := make([]*models.Blog, 0)
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
	sort.SliceStable(blogs, func(i, j int) bool { return blogs[i].Id > blogs[j].Id })
	return blogs, nil
}

func (r *BlogRepositoryOffset) ListByTag(
	ctx context.Context, tx infrastracture.TX, tag string, option *options.ListBlogOptions,
) (models.Blogs, error) {
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
		Order(goqu.I("id").Desc()).
		Select(
			"id", "author_id", "title", "description",
			"thumbnail_image_file_name", "is_public", "created", "modified",
		).
		Limit(uint(option.Limit))
	if option.IsPublic {
		builder = builder.Where(goqu.Ex{"is_public": true})
	}
	if option.Page != nil {
		offset := r.buildOffset(*option.Page, option.Limit)
		builder = builder.Offset(uint(offset))
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
	sort.SliceStable(blogs, func(i, j int) bool { return blogs[i].Id > blogs[j].Id })
	return blogs, nil
}

func (r *BlogRepositoryOffset) ListByKeyword(
	ctx context.Context, tx infrastracture.TX, keyword string, option *options.ListBlogOptions,
) (models.Blogs, error) {
	builder := goqu.
		From("blogs").
		Where(goqu.ExOr{
			"title":       goqu.Op{"like": "%" + keyword + "%"},
			"description": goqu.Op{"like": "%" + keyword + "%"},
		}).
		Order(goqu.I("id").Desc()).
		Select(
			"id", "author_id", "title", "description",
			"thumbnail_image_file_name", "is_public", "created", "modified",
		).
		Limit(uint(option.Limit))
	if option.IsPublic {
		builder = builder.Where(goqu.Ex{"is_public": true})
	}
	if option.Page != nil {
		offset := r.buildOffset(*option.Page, option.Limit)
		builder = builder.Offset(uint(offset))
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
	sort.SliceStable(blogs, func(i, j int) bool { return blogs[i].Id > blogs[j].Id })
	return blogs, nil
}
