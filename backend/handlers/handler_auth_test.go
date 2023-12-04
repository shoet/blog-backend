package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/testutil"
)

func Test_AuthLoginHandler(t *testing.T) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	wantToken := "authtoken"
	tests := []struct {
		name      string
		args      requestBody
		status    int
		want      any
		setCookie bool
	}{
		{
			name: "success",
			args: requestBody{
				Email:    "test@example.com",
				Password: "test",
			},
			status: 200,
			want: struct {
				AuthToken string `json:"authToken"`
			}{
				AuthToken: wantToken,
			},
			setCookie: true,
		},
		{
			name: "validation error",
			args: requestBody{
				Email: "test@example.com",
			},
			status: 400,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "BadRequest",
			},
			setCookie: false,
		},
		{
			name: "unauthorized",
			args: requestBody{
				Email:    "test@example.com",
				Password: "test",
			},
			status: 401,
			want: struct {
				Message string `json:"message"`
			}{
				Message: "Unauthorized",
			},
			setCookie: false,
		},
	}

	validator := validator.New()
	cookie := NewCookieManager("test", "https://test.example.com")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authManagerMock := &AuthManagerMock{}
			authManagerMock.LoginFunc = func(
				ctx context.Context, email string, password string,
			) (string, error) {
				if tt.status == 200 {
					if email != tt.args.Email {
						t.Errorf("email is invalid. got: %s, want: %s", email, tt.args.Email)
					}
					if password != tt.args.Password {
						t.Errorf("password is invalid. got: %s, want: %s", password, tt.args.Password)
					}
					return wantToken, nil
				}
				return "", errors.New("failed to login")
			}

			sut := NewAuthLoginHandler(authManagerMock, validator, cookie)

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

			parser := &http.Request{Header: http.Header{"Cookie": resp.Header["Set-Cookie"]}}
			if tt.setCookie {
				gotCookie, err := parser.Cookie("authToken")
				if errors.Is(err, http.ErrNoCookie) {
					t.Errorf("cookie is not set")
				}
				if gotCookie.Value != wantToken {
					t.Errorf("cookie is invalid. got: %s, want: %s", gotCookie.Value, wantToken)
				}
			} else {
				_, err := parser.Cookie("authToken")
				if !errors.Is(err, http.ErrNoCookie) {
					t.Errorf("cookie is set")
				}
			}
		})
	}

}
