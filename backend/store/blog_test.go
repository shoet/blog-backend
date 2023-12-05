package store

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/testutil"
)

func Test_BlogRepository_Add(t *testing.T) {

	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBMySQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	sut := NewBlogRepository(clocker)

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
					Created:                clocker.Now(),
					Modified:               clocker.Now(),
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
					Created:                clocker.Now(),
					Modified:               clocker.Now(),
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

			row := tx.QueryRowContext(ctx, "SELECT * FROM blogs WHERE id = ?", blogId)
			var got models.Blog
			if err := row.Scan(
				&got.Id, &got.AuthorId, &got.Title, &got.Content, &got.Description,
				&got.ThumbnailImageFileName, &got.IsPublic, &got.Created, &got.Modified,
			); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			opt := cmpopts.IgnoreFields(models.Blog{}, "Id")
			if diff := cmp.Diff(tt.want.blog, &got, opt); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}
