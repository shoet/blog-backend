package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/shoet/blog/testutil"
)

func Test_HealthCheckHandler(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}
	tests := []struct {
		name   string
		status int
		want   response
	}{
		{
			name:   "success",
			status: 200,
			want: response{
				Message: "OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "", nil)

			sut := &HealthCheckHandler{}
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
