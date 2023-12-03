package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
	"github.com/shoet/blog/testutil"
)

func Test_TagListHandler(t *testing.T) {

	type args struct{}

	tests := []struct {
		name    string
		args    args
		status  int
		want    interface{}
		wantErr error
	}{
		{
			name:   "success",
			args:   args{},
			status: 200,
			want: []*models.Tag{
				{
					Id:   1,
					Name: "tag1",
				},
				{
					Id:   2,
					Name: "tag2",
				},
			},
		},
		{
			name:   "internal server error",
			args:   args{},
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
			blogServiceMock.ListTagsFunc = func(
				ctx context.Context, option options.ListTagsOptions,
			) ([]*models.Tag, error) {
				if tt.status == 200 {
					return tt.want.([]*models.Tag), nil
				} else {
					return nil, errors.New("error")
				}
			}

			sut := NewTagListHandler(blogServiceMock)
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
