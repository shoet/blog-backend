package repository_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/internal/testutil"
)

// TODO: ok
func Test_BlogRepository_Add(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepository(clocker)

	type args struct {
		blog *models.Blog
	}

	type want struct {
		blog *models.Blog
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
			},
			want: want{
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()
			blogId, err := sut.Add(ctx, tx, tt.args.blog)
			if err != nil {
				t.Fatalf("failed to add blog: %v", err)
			}

			row := tx.QueryRowContext(ctx, "SELECT * FROM blogs WHERE id = $1", blogId)
			var got models.Blog
			if err := row.Scan(
				&got.Id, &got.AuthorId, &got.Title, &got.Content, &got.Description,
				&got.ThumbnailImageFileName, &got.IsPublic, &got.Created, &got.Modified,
			); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			opt := cmpopts.IgnoreFields(models.Blog{}, "Id", "Created", "Modified")
			if diff := cmp.Diff(tt.want.blog, &got, opt); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func generateTestBlogs(t *testing.T, count int, now time.Time) []*models.Blog {
	t.Helper()
	blogs := make([]*models.Blog, count)
	for i := 0; i < count; i++ {
		b := &models.Blog{
			AuthorId:               1,
			Title:                  fmt.Sprintf("title%d", i),
			Content:                fmt.Sprintf("content%d", i),
			Description:            fmt.Sprintf("description%d", i),
			ThumbnailImageFileName: fmt.Sprintf("thumbnail_image_file_name%d", i),
			IsPublic:               true,
			Created:                uint(now.Unix()),
			Modified:               uint(now.Unix()),
		}
		blogs[i] = b
	}
	return blogs
}

func generateTestBlogsWithPublic(t *testing.T, count int, now time.Time) []*models.Blog {
	t.Helper()
	blogs := make([]*models.Blog, count)
	for i := 0; i < count; i++ {
		b := &models.Blog{
			AuthorId:               1,
			Title:                  fmt.Sprintf("title%d", i),
			Content:                fmt.Sprintf("content%d", i),
			Description:            fmt.Sprintf("description%d", i),
			ThumbnailImageFileName: fmt.Sprintf("thumbnail_image_file_name%d", i),
			IsPublic:               true,
			Created:                uint(now.Unix()),
			Modified:               uint(now.Unix()),
		}
		if i%2 == 0 {
			b.IsPublic = false
		}
		blogs[i] = b
	}
	return blogs
}

func Test_BlogRepository_List(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepository(clocker)

	type args struct {
		limit    *int64
		isPublic bool
		prepare  func(ctx context.Context, tx infrastracture.TX) error
	}

	type want struct {
		count int
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				isPublic: true,
				limit:    func() *int64 { var v int64 = 20; return &v }(),
				prepare: func(ctx context.Context, tx infrastracture.TX) error {
					// blogを20個作成する
					blogs := generateTestBlogs(t, 20, clocker.Now())
					for _, b := range blogs {
						prepareTask := `
						INSERT INTO blogs
							(
								author_id, title, content, description, 
								thumbnail_image_file_name, is_public)
						VALUES
							($1, $2, $3, $4, $5, $6)
						`
						_, err = tx.ExecContext(
							ctx, prepareTask,
							b.AuthorId, b.Title, b.Content, b.Description,
							b.ThumbnailImageFileName, b.IsPublic)
						if err != nil {
							return fmt.Errorf("failed to prepare task: %v", err)
						}

					}
					return nil
				},
			},
			want: want{
				count: 20,
			},
		},
		{
			id: "limit 10",
			args: args{
				limit: func() *int64 { var v int64 = 10; return &v }(),
				prepare: func(ctx context.Context, tx infrastracture.TX) error {
					// blogを20個作成する
					blogs := generateTestBlogs(t, 20, clocker.Now())
					for _, b := range blogs {
						prepareTask := `
						INSERT INTO blogs
							(
								author_id, title, content, description, 
								thumbnail_image_file_name, is_public)
						VALUES
							($1, $2, $3, $4, $5, $6)
						`
						_, err := tx.ExecContext(
							ctx, prepareTask,
							b.AuthorId, b.Title, b.Content, b.Description,
							b.ThumbnailImageFileName, b.IsPublic)
						if err != nil {
							return fmt.Errorf("failed to prepare task: %v", err)
						}

					}
					return nil
				},
			},
			want: want{
				count: 10,
			},
		},
		{
			id: "isNotPublic",
			args: args{
				limit:    nil,
				isPublic: false,
				prepare: func(ctx context.Context, tx infrastracture.TX) error {
					// publicなblogを10個、privateなblogを10個作成する
					blogs := generateTestBlogsWithPublic(t, 20, clocker.Now())
					for _, b := range blogs {
						prepareTask := `
						INSERT INTO blogs
							(
								author_id, title, content, description, 
								thumbnail_image_file_name, is_public)
						VALUES
							($1, $2, $3, $4, $5, $6)
						`
						_, err := tx.ExecContext(
							ctx, prepareTask,
							b.AuthorId, b.Title, b.Content, b.Description,
							b.ThumbnailImageFileName, b.IsPublic)
						if err != nil {
							return fmt.Errorf("failed to prepare task: %v", err)
						}

					}
					return nil
				},
			},
			want: want{
				count: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()

			if err := tt.args.prepare(ctx, tx); err != nil {
				t.Fatalf("failed to prepare: %v", err)
			}

			listOption := options.NewListBlogOptions(&tt.args.isPublic, tt.args.limit)

			blogs, err := sut.List(ctx, tx, listOption)
			if err != nil {
				t.Fatalf("failed to list blogs: %v", err)
			}
			if tt.want.count != len(blogs) {
				t.Fatalf("failed to count blogs: %v", err)
			}

		})
	}
}

// TODO: ok
func Test_BlogRepository_Delete(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepository(clocker)

	type args struct {
		prepareCreateBlog []*models.Blog
	}

	type want struct {
		count int
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success delete 1",
			args: args{
				prepareCreateBlog: generateTestBlogs(t, 10, clocker.Now()),
			},
			want: want{
				count: 9,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()

			for _, b := range tt.args.prepareCreateBlog {
				prepareTask := `
				INSERT INTO blogs
					(
						author_id, title, content, description, 
						thumbnail_image_file_name, is_public)
				VALUES
					($1, $2, $3, $4, $5, $6)
				`
				_, err := tx.ExecContext(
					ctx, prepareTask,
					b.AuthorId, b.Title, b.Content, b.Description,
					b.ThumbnailImageFileName, b.IsPublic)
				if err != nil {
					t.Fatalf("failed to prepare task: %v", err)
				}

			}

			selectSQL := `SELECT id FROM blogs LIMIT 1;`
			row := tx.QueryRowxContext(ctx, selectSQL)
			var blog models.Blog
			if err := row.Scan(&blog.Id); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			if err := sut.Delete(ctx, tx, blog.Id); err != nil {
				t.Fatalf("failed to delete blog: %v", err)
			}
		})
	}

}

// TODO: ok
func Test_BlogRepository_Get(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepository(clocker)

	type args struct {
		prepareCreateBlog []*models.Blog
	}

	type want struct {
		blog *models.Blog
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				prepareCreateBlog: []*models.Blog{
					{
						AuthorId:               1,
						Title:                  "title",
						Content:                "content",
						Description:            "description",
						ThumbnailImageFileName: "thumbnail_image_file_name",
						IsPublic:               true,
					},
				},
			},
			want: want{
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
			},
		},
		{
			id: "failed_not_found",
			args: args{
				prepareCreateBlog: []*models.Blog{},
			},
			want: want{
				blog: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()

			var blogId models.BlogId
			for _, b := range tt.args.prepareCreateBlog {
				prepareTask := `
				INSERT INTO blogs
					(
						author_id, title, content, description, 
						thumbnail_image_file_name, is_public)
				VALUES
					($1, $2, $3, $4, $5, $6)
				RETURNING 
					id
				`
				row := tx.QueryRowxContext(
					ctx, prepareTask,
					b.AuthorId, b.Title, b.Content, b.Description,
					b.ThumbnailImageFileName, b.IsPublic)
				if row.Err() != nil {
					t.Fatalf("failed to prepare task: %v", row.Err())
				}
				var id models.BlogId
				err := row.Scan(&id)
				if err != nil {
					t.Fatalf("failed get Scan: %v", err)
				}
				blogId = models.BlogId(id)
			}

			got, err := sut.Get(ctx, tx, blogId)
			if err != nil {
				t.Fatalf("failed to get blog: %v", err)
			}

			cmpOptions := cmpopts.IgnoreFields(models.Blog{}, "Id", "Created", "Modified")
			if diff := cmp.Diff(tt.want.blog, got, cmpOptions); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}

}

// TODO: ok
func Test_BlogRepository_Put(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepository(clocker)

	type args struct {
		prepareCreateBlog []*models.Blog
		blog              *models.Blog
	}

	type want struct {
		blog *models.Blog
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				prepareCreateBlog: []*models.Blog{
					{
						AuthorId:               1,
						Title:                  "title",
						Content:                "content",
						Description:            "description",
						ThumbnailImageFileName: "thumbnail_image_file_name",
						IsPublic:               true,
					},
				},
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "titleeee",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
			},
			want: want{
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "titleeee",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()

			var blogId models.BlogId
			for _, b := range tt.args.prepareCreateBlog {
				prepareTask := `
				INSERT INTO blogs
					(
						author_id, title, content, description, 
						thumbnail_image_file_name, is_public)
				VALUES
					($1, $2, $3, $4, $5, $6)
				RETURNING
					id
				`
				row := tx.QueryRowxContext(
					ctx, prepareTask,
					b.AuthorId, b.Title, b.Content, b.Description,
					b.ThumbnailImageFileName, b.IsPublic)
				if row.Err() != nil {
					t.Fatalf("failed to prepare task: %v", row.Err())
				}

				var id int64
				err := row.Scan(&id)
				if err != nil {
					t.Fatalf("failed to get last insert id: %v", err)
				}
				blogId = models.BlogId(id)
			}

			tt.args.blog.Id = blogId

			blogId, err := sut.Put(ctx, tx, tt.args.blog)
			if err != nil {
				t.Fatalf("failed to put blog: %v", err)
			}

			selectQuery := `SELECT * FROM blogs WHERE id = $1`
			var got []*models.Blog
			if err := tx.SelectContext(ctx, &got, selectQuery, blogId); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			cmpOptions := cmpopts.IgnoreFields(models.Blog{}, "Id", "Created", "Modified")
			if diff := cmp.Diff(tt.want.blog, got[0], cmpOptions); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}

}

// TODO: ok
func Test_BlogRepository_AddBlogTag(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepository(clocker)

	type blogTag struct {
		Id     int64         `db:"id"`
		BlogId models.BlogId `db:"blog_id"`
		TagId  models.TagId  `db:"tag_id"`
	}

	type args struct {
		prepareBlogTag []*blogTag
		blogTag        *blogTag
	}

	type want struct {
		blogTag []*blogTag
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				prepareBlogTag: []*blogTag{},
				blogTag: &blogTag{
					BlogId: 1,
					TagId:  1,
				},
			},
			want: want{
				blogTag: []*blogTag{
					{
						BlogId: 1,
						TagId:  1,
					},
				},
			},
		},
		{
			id: "success_already_exists",
			args: args{
				prepareBlogTag: []*blogTag{
					{
						BlogId: 1,
						TagId:  1,
					},
					{
						BlogId: 2,
						TagId:  2,
					},
				},
				blogTag: &blogTag{
					BlogId: 1,
					TagId:  1,
				},
			},
			want: want{
				blogTag: []*blogTag{
					{
						BlogId: 1,
						TagId:  1,
					},
					{
						BlogId: 2,
						TagId:  2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()

			for _, bt := range tt.args.prepareBlogTag {
				prepareTask := `
				INSERT INTO blogs_tags
					(blog_id, tag_id)
				VALUES
					($1, $2)
				`
				_, err := tx.ExecContext(ctx, prepareTask, bt.BlogId, bt.TagId)
				if err != nil {
					t.Fatalf("failed to prepare task: %v", err)
				}
			}

			_, err := sut.AddBlogTag(ctx, tx, tt.args.blogTag.BlogId, tt.args.blogTag.TagId)
			if err != nil {
				t.Fatalf("failed to put blog: %v", err)
			}

			selectQuery := `SELECT blog_id, tag_id FROM blogs_tags`
			var got []*blogTag
			if err := tx.SelectContext(ctx, &got, selectQuery); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].BlogId < got[j].BlogId
			})
			sort.Slice(tt.want.blogTag, func(i, j int) bool {
				return tt.want.blogTag[i].BlogId < tt.want.blogTag[j].BlogId
			})

			cmpOption := cmpopts.IgnoreFields(blogTag{}, "Id")
			if diff := cmp.Diff(tt.want.blogTag, got, cmpOption); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}

}

// TODO
func Test_BlogRepository_SelectBlogsTagsByOtherUsingBlog(t *testing.T) {}
func Test_BlogRepository_SelectBlogsTags(t *testing.T)                 {}
func Test_BlogRepository_DeleteBlogsTags(t *testing.T)                 {}
func Test_BlogRepository_SelectTags(t *testing.T)                      {}
func Test_BlogRepository_AddTag(t *testing.T)                          {}
func Test_BlogRepository_DeleteTag(t *testing.T)                       {}
func Test_BlogRepository_ListTags(t *testing.T)                        {}
