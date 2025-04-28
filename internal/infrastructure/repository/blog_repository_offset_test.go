package repository_test

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/infrastructure/repository"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/internal/testutil"
)

func Test_BlogRepositoryOffset_List(t *testing.T) {
	ctx := context.Background()
	clocker := &clocker.FiexedClocker{}
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepositoryOffset(clocker)

	testdata := []*models.Blog{
		{Id: 1, AuthorId: 1, Title: "title1", Content: "content1", Description: "description1", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
		{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
		{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
		{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
		{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
		{Id: 6, AuthorId: 1, Title: "title6", Content: "content6", Description: "description6", ThumbnailImageFileName: "thumbnail6", IsPublic: true},
		{Id: 7, AuthorId: 1, Title: "title7", Content: "content7", Description: "description7", ThumbnailImageFileName: "thumbnail7", IsPublic: true},
		{Id: 8, AuthorId: 1, Title: "title8", Content: "content8", Description: "description8", ThumbnailImageFileName: "thumbnail8", IsPublic: true},
		{Id: 9, AuthorId: 1, Title: "title9", Content: "content9", Description: "description9", ThumbnailImageFileName: "thumbnail9", IsPublic: true},
		{Id: 10, AuthorId: 1, Title: "title10", Content: "content10", Description: "description10", ThumbnailImageFileName: "thumbnail10", IsPublic: true},
		{Id: 11, AuthorId: 1, Title: "title11", Content: "content11", Description: "description11", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
		{Id: 12, AuthorId: 1, Title: "title12", Content: "content12", Description: "description12", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
		{Id: 13, AuthorId: 1, Title: "title13", Content: "content13", Description: "description13", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
		{Id: 14, AuthorId: 1, Title: "title14", Content: "content14", Description: "description14", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
		{Id: 15, AuthorId: 1, Title: "title15", Content: "content15", Description: "description15", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
	}

	type args struct {
		option *options.ListBlogOptions
	}
	type wants struct {
		blogs models.Blogs
		err   error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "単純なLimit",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     1,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 15, AuthorId: 1, Title: "title15", Content: "content15", Description: "description15", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
					{Id: 14, AuthorId: 1, Title: "title14", Content: "content14", Description: "description14", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
					{Id: 13, AuthorId: 1, Title: "title13", Content: "content13", Description: "description13", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
					{Id: 12, AuthorId: 1, Title: "title12", Content: "content12", Description: "description12", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
					{Id: 11, AuthorId: 1, Title: "title11", Content: "content11", Description: "description11", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "ページの指定",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     2,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 10, AuthorId: 1, Title: "title10", Content: "content10", Description: "description10", ThumbnailImageFileName: "thumbnail10", IsPublic: true},
					{Id: 9, AuthorId: 1, Title: "title9", Content: "content9", Description: "description9", ThumbnailImageFileName: "thumbnail9", IsPublic: true},
					{Id: 8, AuthorId: 1, Title: "title8", Content: "content8", Description: "description8", ThumbnailImageFileName: "thumbnail8", IsPublic: true},
					{Id: 7, AuthorId: 1, Title: "title7", Content: "content7", Description: "description7", ThumbnailImageFileName: "thumbnail7", IsPublic: true},
					{Id: 6, AuthorId: 1, Title: "title6", Content: "content6", Description: "description6", ThumbnailImageFileName: "thumbnail6", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "末尾",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    6,
					Page:     3,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "範囲外の指定",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     4,
				},
			},
			wants: wants{
				blogs: []*models.Blog{},
				err:   nil,
			},
		},
		{
			name: "is_public=false",
			args: args{
				option: &options.ListBlogOptions{
					Limit: 5,
					Page:  3,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
					{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
					{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
					{Id: 1, AuthorId: 1, Title: "title1", Content: "content1", Description: "description1", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
				},
				err: nil,
			},
		},
	}

	testdataVals := []goqu.Record{}
	for _, d := range testdata {
		goquVal := goqu.Record{
			"id":                        d.Id,
			"author_id":                 d.AuthorId,
			"title":                     d.Title,
			"content":                   d.Content,
			"description":               d.Description,
			"thumbnail_image_file_name": d.ThumbnailImageFileName,
			"is_public":                 d.IsPublic,
		}
		testdataVals = append(testdataVals, goquVal)
	}

	sql, params, err := goqu.
		Insert("blogs").
		Rows(testdataVals).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := sut.List(ctx, db, tt.args.option)
			if diff := cmp.Diff(tt.wants.err, err); diff != "" {
				t.Errorf("unexpected error: %v", diff)
			}

			options := cmp.Options{
				cmpopts.IgnoreFields(models.Blog{}, "Created", "Modified", "Tags", "Content"),
			}

			if diff := cmp.Diff(tt.wants.blogs, got, options); diff != "" {
				t.Errorf("unexpected blogs: %v", diff)
			}

		})
	}

	sql, params, err = goqu.
		Delete("blogs").
		Where(goqu.I("id").Gte(1)).
		Where(goqu.I("id").Lte(15)).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to delete test data: %v", err)
	}
}

func Test_BlogRepositoryOffset_ListByTag(t *testing.T) {
	ctx := context.Background()
	clocker := &clocker.FiexedClocker{}
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepositoryOffset(clocker)

	testdataTags := []*models.Tag{
		{Id: 1, Name: "tag1"},
		{Id: 2, Name: "tag2"},
		{Id: 3, Name: "tag3"},
	}

	testdataBlog := []*models.Blog{
		{Id: 1, AuthorId: 1, Title: "title1", Content: "content1", Description: "description1", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
		{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
		{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
		{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
		{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
		{Id: 6, AuthorId: 1, Title: "title6", Content: "content6", Description: "description6", ThumbnailImageFileName: "thumbnail6", IsPublic: true},
		{Id: 7, AuthorId: 1, Title: "title7", Content: "content7", Description: "description7", ThumbnailImageFileName: "thumbnail7", IsPublic: true},
		{Id: 8, AuthorId: 1, Title: "title8", Content: "content8", Description: "description8", ThumbnailImageFileName: "thumbnail8", IsPublic: true},
		{Id: 9, AuthorId: 1, Title: "title9", Content: "content9", Description: "description9", ThumbnailImageFileName: "thumbnail9", IsPublic: true},
		{Id: 10, AuthorId: 1, Title: "title10", Content: "content10", Description: "description10", ThumbnailImageFileName: "thumbnail10", IsPublic: true},
		{Id: 11, AuthorId: 1, Title: "title11", Content: "content11", Description: "description11", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
		{Id: 12, AuthorId: 1, Title: "title12", Content: "content12", Description: "description12", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
		{Id: 13, AuthorId: 1, Title: "title13", Content: "content13", Description: "description13", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
		{Id: 14, AuthorId: 1, Title: "title14", Content: "content14", Description: "description14", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
		{Id: 15, AuthorId: 1, Title: "title15", Content: "content15", Description: "description15", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
	}

	testdataBlogTags := []*models.BlogsTags{
		{BlogId: 15, TagId: 1},
		{BlogId: 14, TagId: 1},
		{BlogId: 13, TagId: 1},
		{BlogId: 12, TagId: 1},
		{BlogId: 11, TagId: 1},
		{BlogId: 5, TagId: 3},
		{BlogId: 4, TagId: 3},
		{BlogId: 3, TagId: 3},
		{BlogId: 2, TagId: 3},
		{BlogId: 1, TagId: 3},
	}

	type args struct {
		option *options.ListBlogOptions
		tag    string
	}
	type wants struct {
		blogs models.Blogs
		err   error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "単純なLimit",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    3,
					Page:     1,
				},
				tag: "tag1",
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 15, AuthorId: 1, Title: "title15", Content: "content15", Description: "description15", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
					{Id: 14, AuthorId: 1, Title: "title14", Content: "content14", Description: "description14", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
					{Id: 13, AuthorId: 1, Title: "title13", Content: "content13", Description: "description13", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "ページの指定",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    3,
					Page:     2,
				},
				tag: "tag1",
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 12, AuthorId: 1, Title: "title12", Content: "content12", Description: "description12", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
					{Id: 11, AuthorId: 1, Title: "title11", Content: "content11", Description: "description11", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "is_public=true",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     1,
				},
				tag: "tag3",
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
					{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
					{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			// is_public=falseが結果に入ってくる
			name: "is_public=false",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: false,
					Limit:    5,
					Page:     1,
				},
				tag: "tag3",
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
					{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
					{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
					{Id: 1, AuthorId: 1, Title: "title1", Content: "content1", Description: "description1", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
				},
				err: nil,
			},
		},
	}

	func(db *sqlx.DB) {
		testdataVals := []goqu.Record{}
		for _, d := range testdataBlog {
			goquVal := goqu.Record{
				"id":                        d.Id,
				"author_id":                 d.AuthorId,
				"title":                     d.Title,
				"content":                   d.Content,
				"description":               d.Description,
				"thumbnail_image_file_name": d.ThumbnailImageFileName,
				"is_public":                 d.IsPublic,
			}
			testdataVals = append(testdataVals, goquVal)
		}

		sql, params, err := goqu.
			Insert("blogs").
			Rows(testdataVals).
			ToSQL()
		if err != nil {
			t.Fatalf("failed to build sql: %v", err)
		}
		_, err = db.ExecContext(ctx, sql, params...)
		if err != nil {
			t.Fatalf("failed to insert test data blogs: %v", err)
		}
	}(db)

	func(db *sqlx.DB) {
		testdataVals := []goqu.Record{}
		for _, d := range testdataTags {
			goquVal := goqu.Record{
				"id":   d.Id,
				"name": d.Name,
			}
			testdataVals = append(testdataVals, goquVal)
		}

		sql, params, err := goqu.
			Insert("tags").
			Rows(testdataVals).
			ToSQL()
		if err != nil {
			t.Fatalf("failed to build sql: %v", err)
		}
		_, err = db.ExecContext(ctx, sql, params...)
		if err != nil {
			t.Fatalf("failed to insert test data tags: %v", err)
		}
	}(db)

	func(db *sqlx.DB) {
		testdataVals := []goqu.Record{}
		for _, d := range testdataBlogTags {
			goquVal := goqu.Record{
				"blog_id": d.BlogId,
				"tag_id":  d.TagId,
			}
			testdataVals = append(testdataVals, goquVal)
		}

		sql, params, err := goqu.
			Insert("blogs_tags").
			Rows(testdataVals).
			ToSQL()
		if err != nil {
			t.Fatalf("failed to build sql: %v", err)
		}
		_, err = db.ExecContext(ctx, sql, params...)
		if err != nil {
			t.Fatalf("failed to insert test data blogs_tags: %v", err)
		}
	}(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := sut.ListByTag(ctx, db, tt.args.tag, tt.args.option)
			if diff := cmp.Diff(tt.wants.err, err); diff != "" {
				t.Errorf("unexpected error: %v", diff)
			}

			options := cmp.Options{
				cmpopts.IgnoreFields(models.Blog{}, "Created", "Modified", "Content"),
			}

			if diff := cmp.Diff(tt.wants.blogs, got, options); diff != "" {
				t.Errorf("unexpected blogs: %v", diff)
			}

		})
	}

	sql, params, err := goqu.
		Delete("blogs_tags").
		Where(goqu.I("blog_id").Gte(1)).
		Where(goqu.I("blog_id").Lte(15)).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to delete test data blogs_tags: %v", err)
	}

	sql, params, err = goqu.
		Delete("tags").
		Where(goqu.I("id").Gte(1)).
		Where(goqu.I("id").Lte(5)).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to delete test data: %v", err)
	}

	sql, params, err = goqu.
		Delete("blogs").
		Where(goqu.I("id").Gte(1)).
		Where(goqu.I("id").Lte(15)).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to delete test data: %v", err)
	}

}

func Test_BlogRepositoryOffset_ListByKeyword(t *testing.T) {
	ctx := context.Background()
	clocker := &clocker.FiexedClocker{}
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepositoryOffset(clocker)

	testdata := []*models.Blog{
		{Id: 1, AuthorId: 1, Title: "title1XXX", Content: "content1AAA", Description: "description1AAA", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
		{Id: 2, AuthorId: 1, Title: "title2XXX", Content: "content2AAA", Description: "description2AAA", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
		{Id: 3, AuthorId: 1, Title: "title3XXX", Content: "content3AAA", Description: "description3AAA", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
		{Id: 4, AuthorId: 1, Title: "title4XXX", Content: "content4AAA", Description: "description4AAA", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
		{Id: 5, AuthorId: 1, Title: "title5XXX", Content: "content5AAA", Description: "description5AAA", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
		{Id: 6, AuthorId: 1, Title: "title6YYY", Content: "content6BBB", Description: "description6BBB", ThumbnailImageFileName: "thumbnail6", IsPublic: true},
		{Id: 7, AuthorId: 1, Title: "title7YYY", Content: "content7BBB", Description: "description7BBB", ThumbnailImageFileName: "thumbnail7", IsPublic: true},
		{Id: 8, AuthorId: 1, Title: "title8YYY", Content: "content8BBB", Description: "description8BBB", ThumbnailImageFileName: "thumbnail8", IsPublic: true},
		{Id: 9, AuthorId: 1, Title: "title9YYY", Content: "content9BBB", Description: "description9BBB", ThumbnailImageFileName: "thumbnail9", IsPublic: true},
		{Id: 10, AuthorId: 1, Title: "title10YYY", Content: "content10BBB", Description: "description10BBB", ThumbnailImageFileName: "thumbnail10", IsPublic: true},
		{Id: 11, AuthorId: 1, Title: "title11ZZZ", Content: "content11CCC", Description: "description11CCC", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
		{Id: 12, AuthorId: 1, Title: "title12ZZZ", Content: "content12CCC", Description: "description12CCC", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
		{Id: 13, AuthorId: 1, Title: "title13ZZZ", Content: "content13CCC", Description: "description13CCC", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
		{Id: 14, AuthorId: 1, Title: "title14ZZZ", Content: "content14CCC", Description: "description14CCC", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
		{Id: 15, AuthorId: 1, Title: "title15ZZZ", Content: "content15CCC", Description: "description15CCC", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
	}

	type args struct {
		keyword string
		option  *options.ListBlogOptions
	}
	type wants struct {
		blogs models.Blogs
		err   error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "Descriptionで絞り込んで2ページ目を取得する",
			args: args{
				keyword: "AAA",
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    2,
					Page:     2,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 3, AuthorId: 1, Title: "title3XXX", Content: "content3AAA", Description: "description3AAA", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2XXX", Content: "content2AAA", Description: "description2AAA", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "Titleで絞り込んで2ページ目を取得する",
			args: args{
				keyword: "ZZZ",
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    3,
					Page:     2,
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 12, AuthorId: 1, Title: "title12ZZZ", Content: "content12CCC", Description: "description12CCC", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
					{Id: 11, AuthorId: 1, Title: "title11ZZZ", Content: "content11CCC", Description: "description11CCC", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
				},
				err: nil,
			},
		},
	}

	testdataVals := []goqu.Record{}
	for _, d := range testdata {
		goquVal := goqu.Record{
			"id":                        d.Id,
			"author_id":                 d.AuthorId,
			"title":                     d.Title,
			"content":                   d.Content,
			"description":               d.Description,
			"thumbnail_image_file_name": d.ThumbnailImageFileName,
			"is_public":                 d.IsPublic,
		}
		testdataVals = append(testdataVals, goquVal)
	}

	sql, params, err := goqu.
		Insert("blogs").
		Rows(testdataVals).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := sut.ListByKeyword(ctx, db, tt.args.keyword, tt.args.option)
			if diff := cmp.Diff(tt.wants.err, err); diff != "" {
				t.Errorf("unexpected error: %v", diff)
			}

			options := cmp.Options{
				cmpopts.IgnoreFields(models.Blog{}, "Created", "Modified", "Tags", "Content"),
			}

			if diff := cmp.Diff(tt.wants.blogs, got, options); diff != "" {
				t.Errorf("unexpected blogs: %v", diff)
			}

		})
	}

	sql, params, err = goqu.
		Delete("blogs").
		Where(goqu.I("id").Gte(1)).
		Where(goqu.I("id").Lte(15)).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to delete test data: %v", err)
	}
}
