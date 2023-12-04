package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
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
				if tt.status == 200 {
					return tt.want.([]*models.Blog), nil
				}
				return nil, errors.New("internal server error")
			}

			sut := NewBlogListHandler(blogServiceMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/blogs", nil)
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
			r := httptest.NewRequest("GET", "/blogs", nil)
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
			r := httptest.NewRequest("GET", "/blogs", nil)
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
