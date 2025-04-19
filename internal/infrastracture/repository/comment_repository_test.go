package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_CommentRepository_CreateComment(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewCommentRepository(clocker)
	_ = sut

	type args struct {
		blog     *models.Blog
		userId   *models.UserId
		clientId *string
		content  string
	}

	type want struct {
		Comment *models.Comment
		err     error
	}

	strPtr := func(s string) *string {
		return &s
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
					Id:                     1,
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
				userId:   nil,
				clientId: strPtr("abcd"),
				content:  "comment",
			},
			want: want{
				Comment: &models.Comment{
					CommentId: 1,
					BlogId:    1,
					UserId:    nil,
					ClientId:  strPtr("abcd"),
					Content:   "comment",
					IsEdited:  false,
					IsDeleted: false,
					Created:   clocker.Now().Unix(),
					Modified:  clocker.Now().Unix(),
				},
				err: nil,
			},
		},
		{
			id: "ブログ記事が存在しない場合はエラー",
			args: args{
				blog: &models.Blog{
					Id:                     1,
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
				userId:   nil,
				clientId: strPtr("abcd"),
				content:  "comment",
			},
			want: want{
				Comment: &models.Comment{
					CommentId: 1,
					BlogId:    2,
					UserId:    nil,
					ClientId:  strPtr("abcd"),
					Content:   "comment",
					IsEdited:  false,
					IsDeleted: false,
					Created:   clocker.Now().Unix(),
					Modified:  clocker.Now().Unix(),
				},
				err: errors.New("blog not found"),
			},
		},
		{
			id: "user_idとclient_idの両方がnilの場合はエラー",
			args: args{
				blog: &models.Blog{
					Id:                     1,
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
				},
				userId:   nil,
				clientId: strPtr("abcd"),
				content:  "comment",
			},
			want: want{
				Comment: &models.Comment{
					CommentId: 1,
					BlogId:    2,
					UserId:    nil,
					ClientId:  nil,
					Content:   "comment",
					IsEdited:  false,
					IsDeleted: false,
					Created:   clocker.Now().Unix(),
					Modified:  clocker.Now().Unix(),
				},
				err: errors.New("user_id and client_id cannot be both nil"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()
			insert := `
			INSERT INTO blogs (
				id,
				author_id, title, content, description,
				thumbnail_image_file_name, is_public, created, modified
			)
			VALUES (
				$1, 
				$2, $3, $4,
				$5, $6, $7, $8, $9
			)
			`
			_, err := tx.ExecContext(ctx, insert,
				tt.args.blog.Id,
				tt.args.blog.AuthorId, tt.args.blog.Title, tt.args.blog.Content,
				tt.args.blog.Description, tt.args.blog.ThumbnailImageFileName,
				tt.args.blog.IsPublic,
				clocker.Now().Unix(), clocker.Now().Unix(),
			)
			if err != nil {
				t.Fatalf("failed to insert blog: %v", err)
			}

			commentId, gotErr := sut.CreateComment(
				ctx, tx,
				tt.args.blog.Id, tt.args.userId, tt.args.clientId, tt.args.content,
			)
			if gotErr != nil {
				if tt.want.err != nil {
					assert.Error(t, gotErr)
				} else {
					t.Fatalf("failed to create comment: %v", gotErr)
				}
			}

			row := tx.QueryRowContext(ctx, "SELECT * FROM comments WHERE comment_id = $1", commentId)
			var got models.Comment
			if err := row.Scan(
				&got.CommentId, &got.BlogId, &got.ClientId, &got.UserId,
				&got.Content, &got.IsEdited, &got.IsDeleted,
				&got.Created, &got.Modified,
			); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			if tt.want.err == nil {
				opt := cmpopts.IgnoreFields(models.Comment{}, "CommentId")
				want := &models.Comment{
					CommentId: tt.want.Comment.CommentId,
					BlogId:    tt.want.Comment.BlogId,
					ClientId:  tt.want.Comment.ClientId,
					UserId:    tt.want.Comment.UserId,
					Content:   tt.want.Comment.Content,
					IsEdited:  tt.want.Comment.IsEdited,
					IsDeleted: tt.want.Comment.IsDeleted,
					Created:   tt.want.Comment.Created,
					Modified:  tt.want.Comment.Modified,
				}
				if diff := cmp.Diff(want, &got, opt); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			}

		})
	}
}
