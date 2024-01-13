package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/services"
	"github.com/shoet/blog/testutil"
)

func Test_BlogListHandler(t *testing.T) {
	clocker := &clocker.FiexedClocker{}

	tests := []struct {
		name   string
		status int
		want   interface{}
	}{
		{
			name:   "success",
			status: 200,
			want: []*models.Blog{
				{
					Id:       1,
					Title:    "test",
					Content:  "test",
					IsPublic: true,
					Created:  clocker.Now(),
					Modified: clocker.Now(),
				},
			},
		},
		{
			name:   "internal server error",
			status: 500,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "InternalServerError",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.ListBlogFunc = func(
				ctx context.Context, option options.ListBlogOptions,
			) ([]*models.Blog, error) {
				if !option.IsPublic {
					t.Errorf("want: %v, got: %v", true, option.IsPublic)
				}
				if tt.status == 200 {
					return tt.want.([]*models.Blog), nil
				}
				return nil, errors.New("internal server error")
			}

			sut := NewBlogListHandler(blogServiceMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r = testutil.SetLoggerContextToRequest(t, r)

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_BlogGetHandler(t *testing.T) {
	clocker := &clocker.FiexedClocker{}

	type args struct {
		pathParamsBlogId string
	}

	tests := []struct {
		name   string
		args   args
		status int
		want   interface{}
	}{
		{
			name: "success",
			args: args{
				pathParamsBlogId: "1",
			},
			status: 200,
			want: &models.Blog{
				Id:       1,
				Title:    "test",
				Content:  "test",
				IsPublic: true,
				Created:  clocker.Now(),
				Modified: clocker.Now(),
			},
		},
		{
			name: "id is empty",
			args: args{
				pathParamsBlogId: "",
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
		},
		{
			name: "id is invalid format",
			args: args{
				pathParamsBlogId: "a",
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
		},
		{
			name: "entity is not found",
			args: args{
				pathParamsBlogId: "1",
			},
			status: 404,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "NotFound",
			},
		},
		{
			name: "blog service returns error",
			args: args{
				pathParamsBlogId: "1",
			},
			status: 500,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "InternalServerError",
			},
		},
	}

	jwter := &services.JWTerMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.GetBlogFunc = func(
				ctx context.Context, id models.BlogId,
			) (*models.Blog, error) {
				if tt.status == 200 {
					blogId, err := strconv.Atoi(tt.args.pathParamsBlogId)
					if err != nil {
						t.Fatalf("cannot convert blogId to int: %v", err)
					}
					if id != models.BlogId(blogId) {
						t.Errorf("want: %v, got: %v", tt.args.pathParamsBlogId, id)
					}
				}
				if tt.status == 404 {
					return nil, nil
				}
				if tt.status == 500 {
					return nil, errors.New("internal server error")
				}
				return tt.want.(*models.Blog), nil
			}

			sut := NewBlogGetHandler(blogServiceMock, jwter)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.pathParamsBlogId)
			r = testutil.SetLoggerContextToRequest(t, r)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_BlogGetHandlerWithSecret(t *testing.T) {
	clocker := &clocker.FiexedClocker{}

	type args struct {
		pathParamsBlogId string
		authToken        string
	}

	wantAuthToken := "authToken"

	tests := []struct {
		name   string
		args   args
		status int
		want   interface{}
	}{
		{
			name: "get secret item",
			args: args{
				pathParamsBlogId: "1",
			},
			status: 404,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "NotFound",
			},
		},
		{
			name: "get secret item success",
			args: args{
				pathParamsBlogId: "1",
				authToken:        wantAuthToken,
			},
			status: 200,
			want: &models.Blog{
				Id:       1,
				AuthorId: 1,
				Title:    "test",
				Content:  "test",
				IsPublic: false,
				Created:  clocker.Now(),
				Modified: clocker.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.GetBlogFunc = func(
				ctx context.Context, id models.BlogId,
			) (*models.Blog, error) {
				if tt.status == 200 {
					blogId, err := strconv.Atoi(tt.args.pathParamsBlogId)
					if err != nil {
						t.Fatalf("cannot convert blogId to int: %v", err)
					}
					if id != models.BlogId(blogId) {
						t.Errorf("want: %v, got: %v", tt.args.pathParamsBlogId, id)
					}
				}
				if tt.status == 404 {
					return nil, nil
				}
				if tt.status == 500 {
					return nil, errors.New("internal server error")
				}
				return tt.want.(*models.Blog), nil
			}

			jwter := &services.JWTerMock{}
			jwter.VerifyTokenFunc = func(ctx context.Context, token string) (models.UserId, error) {
				if tt.args.authToken != wantAuthToken {
					t.Errorf("want: %v, got: %v", wantAuthToken, tt.args.authToken)
				}
				return tt.want.(*models.Blog).AuthorId, nil
			}

			sut := NewBlogGetHandler(blogServiceMock, jwter)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.pathParamsBlogId)
			r = testutil.SetLoggerContextToRequest(t, r)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			if tt.args.authToken != "" {
				r.Header.Set("Authorization", "Bearer "+tt.args.authToken)
			}

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_BlogAddHandler(t *testing.T) {
	wantBlog := &models.Blog{
		Id:                     1,
		AuthorId:               1,
		Title:                  "test",
		Content:                "test",
		Description:            "test",
		ThumbnailImageFileName: "test",
		IsPublic:               false,
		Tags:                   []string{"test"},
	}

	type requestBody struct {
		Id                     models.BlogId `json:"id" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		Title                  string        `json:"title"`
		Content                string        `json:"content"`
		Description            string        `json:"description"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName"`
		IsPublic               bool          `json:"isPublic"`
		Tags                   []string      `json:"tags"`
	}

	type args struct {
		requestBody requestBody
	}

	tests := []struct {
		id     string
		args   args
		status int
		want   interface{}
	}{
		{

			id: "success_normal",
			args: args{
				requestBody: requestBody{
					AuthorId:               wantBlog.AuthorId,
					Title:                  wantBlog.Title,
					Content:                wantBlog.Content,
					Description:            wantBlog.Description,
					ThumbnailImageFileName: wantBlog.ThumbnailImageFileName,
					IsPublic:               wantBlog.IsPublic,
					Tags:                   wantBlog.Tags,
				},
			},
			status: 200,
			want:   wantBlog,
		},
		{
			id: "failed_validation_error_description",
			args: args{
				requestBody: requestBody{
					AuthorId:               wantBlog.AuthorId,
					Title:                  wantBlog.Title,
					Content:                wantBlog.Content,
					ThumbnailImageFileName: wantBlog.ThumbnailImageFileName,
					IsPublic:               wantBlog.IsPublic,
					Tags:                   wantBlog.Tags,
				},
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
		},
		{
			id: "failed_internal_server_error",
			args: args{
				requestBody: requestBody{
					AuthorId:               wantBlog.AuthorId,
					Title:                  wantBlog.Title,
					Content:                wantBlog.Content,
					Description:            wantBlog.Description,
					ThumbnailImageFileName: wantBlog.ThumbnailImageFileName,
					IsPublic:               wantBlog.IsPublic,
					Tags:                   wantBlog.Tags,
				},
			},
			status: 500,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "InternalServerError",
			},
		},
	}

	validator := validator.New()
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {

			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.AddBlogFunc = func(
				ctx context.Context, blog *models.Blog,
			) (*models.Blog, error) {
				if tt.id == "failed_internal_server_error" {
					return nil, errors.New("internal server error")
				}
				opt := cmpopts.IgnoreFields(models.Blog{}, "Id")
				if diff := cmp.Diff(tt.want, blog, opt); diff != "" {
					t.Errorf("want: %v, got: %v", tt.want, blog)
				}
				return tt.want.(*models.Blog), nil
			}

			sut := NewBlogAddHandler(blogServiceMock, validator)

			var buffer bytes.Buffer
			if err := json.NewEncoder(&buffer).Encode(tt.args.requestBody); err != nil {
				t.Fatalf("failed to encode request body: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", &buffer)
			r = testutil.SetLoggerContextToRequest(t, r)

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}

		})
	}
}

func Test_BlogDeleteHandler(t *testing.T) {
	wantBlogId := models.BlogId(1)

	type requestBody struct {
		Id      models.BlogId `json:"id"`
		IdDummy models.BlogId `json:"idDummy,omitempty"`
	}

	type args struct {
		requestBody requestBody
	}

	tests := []struct {
		id     string
		args   args
		status int
		want   interface{}
	}{
		{
			id: "success_normal",
			args: args{
				requestBody: requestBody{
					Id: wantBlogId,
				},
			},
			status: 200,
			want: struct {
				BlogId models.BlogId `json:"id"`
			}{
				BlogId: wantBlogId,
			},
		},
		{
			id: "failed_validation_error_id",
			args: args{
				requestBody: requestBody{
					IdDummy: wantBlogId,
				},
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
		},
		{
			id: "failed_internal_server_error",
			args: args{
				requestBody: requestBody{
					Id: wantBlogId,
				},
			},
			status: 500,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "InternalServerError",
			},
		},
	}

	validator := validator.New()

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.DeleteBlogFunc = func(ctx context.Context, id models.BlogId) error {
				if tt.id == "success_normal" {
					if id != wantBlogId {
						t.Errorf("want: %v, got: %v", wantBlogId, id)
					}
				}
				if tt.id == "failed_internal_server_error" {
					return errors.New("internal server error")
				}
				return nil
			}

			sut := NewBlogDeleteHandler(blogServiceMock, validator)

			var buffer bytes.Buffer
			if err := json.NewEncoder(&buffer).Encode(tt.args.requestBody); err != nil {
				t.Fatalf("failed to encode request body: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", &buffer)
			r = testutil.SetLoggerContextToRequest(t, r)

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_BlogPutHandler(t *testing.T) {
	wantBlog := &models.Blog{
		Id:                     1,
		AuthorId:               1,
		Title:                  "test",
		Content:                "test",
		Description:            "test",
		ThumbnailImageFileName: "test",
		IsPublic:               false,
		Tags:                   []string{"test"},
	}

	type requestBody struct {
		Id                     models.BlogId `json:"id" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		Title                  string        `json:"title"`
		Content                string        `json:"content"`
		Description            string        `json:"description"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName"`
		IsPublic               bool          `json:"isPublic"`
		Tags                   []string      `json:"tags"`
	}

	type args struct {
		requestBody requestBody
	}

	tests := []struct {
		id     string
		args   args
		status int
		want   interface{}
	}{
		{

			id: "success_normal",
			args: args{
				requestBody: requestBody{
					Id:                     wantBlog.Id,
					AuthorId:               wantBlog.AuthorId,
					Title:                  wantBlog.Title,
					Content:                wantBlog.Content,
					Description:            wantBlog.Description,
					ThumbnailImageFileName: wantBlog.ThumbnailImageFileName,
					IsPublic:               wantBlog.IsPublic,
					Tags:                   wantBlog.Tags,
				},
			},
			status: 200,
			want:   wantBlog,
		},
		{
			id: "failed_validation_error_description",
			args: args{
				requestBody: requestBody{
					Id:                     wantBlog.Id,
					Title:                  wantBlog.Title,
					Content:                wantBlog.Content,
					ThumbnailImageFileName: wantBlog.ThumbnailImageFileName,
					IsPublic:               wantBlog.IsPublic,
					Tags:                   wantBlog.Tags,
				},
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
		},
		{
			id: "failed_internal_server_error",
			args: args{
				requestBody: requestBody{
					Id:                     wantBlog.Id,
					AuthorId:               wantBlog.AuthorId,
					Title:                  wantBlog.Title,
					Content:                wantBlog.Content,
					Description:            wantBlog.Description,
					ThumbnailImageFileName: wantBlog.ThumbnailImageFileName,
					IsPublic:               wantBlog.IsPublic,
					Tags:                   wantBlog.Tags,
				},
			},
			status: 500,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "InternalServerError",
			},
		},
	}

	validator := validator.New()
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {

			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.PutBlogFunc = func(
				ctx context.Context, blog *models.Blog,
			) (*models.Blog, error) {
				if tt.id == "failed_internal_server_error" {
					return nil, errors.New("internal server error")
				}
				if diff := cmp.Diff(tt.want, blog); diff != "" {
					t.Errorf("want: %v, got: %v", tt.want, blog)
				}
				return tt.want.(*models.Blog), nil
			}

			sut := NewBlogPutHandler(blogServiceMock, validator)

			var buffer bytes.Buffer
			if err := json.NewEncoder(&buffer).Encode(tt.args.requestBody); err != nil {
				t.Fatalf("failed to encode request body: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", &buffer)
			r = testutil.SetLoggerContextToRequest(t, r)

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}

		})
	}
}

func Test_BlogListAdminHandler(t *testing.T) {
	clocker := &clocker.FiexedClocker{}

	tests := []struct {
		name   string
		status int
		want   interface{}
	}{
		{
			name:   "success",
			status: 200,
			want: []*models.Blog{
				{
					Id:       1,
					Title:    "test",
					Content:  "test",
					IsPublic: true,
					Created:  clocker.Now(),
					Modified: clocker.Now(),
				},
			},
		},
		{
			name:   "internal server error",
			status: 500,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "InternalServerError",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blogServiceMock := &BlogManagerMock{}
			blogServiceMock.ListBlogFunc = func(
				ctx context.Context, option options.ListBlogOptions,
			) ([]*models.Blog, error) {
				if option.IsPublic {
					t.Errorf("want: %v, got: %v", false, option.IsPublic)
				}
				if tt.status == 200 {
					return tt.want.([]*models.Blog), nil
				}
				return nil, errors.New("internal server error")
			}

			sut := NewBlogListAdminHandler(blogServiceMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r = testutil.SetLoggerContextToRequest(t, r)

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}
		})
	}
}
