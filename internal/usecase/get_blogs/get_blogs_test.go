package get_blogs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/internal/testutil"
)

func Test_GetBlogsUsecase_GetBlogs(t *testing.T) {
	testClocker := clocker.NewFixedClocker()
	testNow := testClocker.Now()
	type args struct {
		prepare         func(ctx context.Context, tx infrastracture.TX) error
		listBlogOptions options.ListBlogOptions
	}
	type wants struct {
		blogs []*models.Blog
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "success public",
			args: args{
				prepare: func(ctx context.Context, tx infrastracture.TX) error {
					query := `
				INSERT INTO blogs
					(
						author_id, title, content, description, 
						thumbnail_image_file_name, is_public, created, modified)
				VALUES
					(?, ?, ?, ?, ?, ?, ?, ?),
					(?, ?, ?, ?, ?, ?, ?, ?)
				`
					if _, err := tx.ExecContext(ctx, query,
						1, "public", "content", "description", "image_url", true, testNow, testNow,
						1, "private", "content", "description", "image_url", false, testNow, testNow,
					); err != nil {
						return fmt.Errorf("failed to prepare for test")
					}
					return nil
				},
				listBlogOptions: options.ListBlogOptions{
					IsPublic: true,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{
						AuthorId:               1,
						Title:                  "public",
						Description:            "description",
						ThumbnailImageFileName: "image_url",
						IsPublic:               true,
						Created:                testNow,
						Modified:               testNow,
					},
				},
			},
		},
		{
			name: "success with private",
			args: args{
				prepare: func(ctx context.Context, tx infrastracture.TX) error {
					query := `
				INSERT INTO blogs
					(
						author_id, title, content, description, 
						thumbnail_image_file_name, is_public, created, modified)
				VALUES
					(?, ?, ?, ?, ?, ?, ?, ?),
					(?, ?, ?, ?, ?, ?, ?, ?)
				`
					if _, err := tx.ExecContext(ctx, query,
						1, "public", "content", "description", "image_url", true, testNow, testNow,
						1, "private", "content", "description", "image_url", false, testNow, testNow,
					); err != nil {
						return fmt.Errorf("failed to prepare for test")
					}
					return nil
				},
				listBlogOptions: options.ListBlogOptions{
					IsPublic: false,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{
						AuthorId:               1,
						Title:                  "public",
						Description:            "description",
						ThumbnailImageFileName: "image_url",
						IsPublic:               true,
						Created:                testNow,
						Modified:               testNow,
					},
					{
						AuthorId:               1,
						Title:                  "private",
						Description:            "description",
						ThumbnailImageFileName: "image_url",
						IsPublic:               false,
						Created:                testNow,
						Modified:               testNow,
					},
				},
			},
		},
	}

	ctx := context.Background()
	db, err := testutil.NewDBMySQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to connect database for test")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blogRepository := repository.NewBlogRepository(testClocker)
			tx, err := db.BeginTxx(ctx, nil)
			if err != nil {
				t.Fatalf("failed to begin transaction")
			}
			if err := tt.args.prepare(ctx, tx); err != nil {
				t.Fatalf("failed to prepare for test")
			}
			blogs, err := blogRepository.List(ctx, tx, tt.args.listBlogOptions)
			if err != nil {
				t.Fatalf("failed to list blogs")
			}
			cmpOpt := cmpopts.IgnoreFields(models.Blog{}, "Id")
			if diff := cmp.Diff(blogs, tt.wants.blogs, cmpOpt); diff != "" {
				t.Errorf("blogs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
