package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/interfaces/handler"
	"github.com/shoet/blog/internal/testutil"
)

func Test_GenerateThumbnailImageSignedURLHandler(t *testing.T) {
	type args struct {
		FileName string
	}
	tests := []struct {
		name   string
		args   interface{}
		status int
		want   interface{}
	}{
		{
			name: "success",
			args: struct {
				FileName string
			}{
				FileName: "test.jpg",
			},
			status: 200,
			want: struct {
				SignedUrl string `json:"signedUrl"`
				PutedUrl  string `json:"putUrl"`
			}{
				SignedUrl: "signed url",
				PutedUrl:  "put url",
			},
		},
		{
			name: "validation error",
			args: struct {
				Name string
			}{
				Name: "test.jpg",
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
		},
		{
			name: "internal server error",
			args: struct {
				FileName string
			}{
				FileName: "test.jpg",
			},
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
			storageServiceMock := &handler.StoragerMock{}
			storageServiceMock.GenerateThumbnailPutURLFunc = func(fileName string) (string, string, error) {
				if tt.status == 500 {
					return "", "", fmt.Errorf("InternalServerError")
				}
				return "signed url", "put url", nil
			}

			sut := handler.NewGenerateThumbnailImageSignedURLHandler(storageServiceMock, validator.New())

			var buffer bytes.Buffer
			if err := json.NewEncoder(&buffer).Encode(tt.args); err != nil {
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
