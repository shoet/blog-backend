package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastracture/models"
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

func Test_AuthSessionLoginHandler(t *testing.T) {
	wantToken := "authtoken"
	wantUser := &models.User{
		Id:    1,
		Name:  "test",
		Email: "test@example.com",
	}
	type args struct {
		authToken           string
		setToken            bool
		headerFormatInvalid bool
	}
	type want struct {
		user      *models.User
		authToken string
		response  interface{}
	}
	tests := []struct {
		name   string
		args   args
		status int
		want   want
	}{
		{
			name: "success",
			args: args{
				authToken: wantToken,
				setToken:  true,
			},
			status: 200,
			want: want{
				user:      wantUser,
				authToken: wantToken,
				response:  wantUser,
			},
		},
		{
			name: "token no set",
			args: args{
				setToken: false,
			},
			status: 401,
			want: want{
				user:      nil,
				authToken: "",
				response: struct {
					Message string `json:"message"`
				}{
					Message: "Unauthorized",
				},
			},
		},
		{
			name: "token format invalid",
			args: args{
				setToken:            true,
				headerFormatInvalid: true,
			},
			status: 401,
			want: want{
				user:      nil,
				authToken: wantToken,
				response: struct {
					Message string `json:"message"`
				}{
					Message: "Unauthorized",
				},
			},
		},
		{
			name: "token is empty",
			args: args{
				setToken:  true,
				authToken: "",
			},
			status: 401,
			want: want{
				user:      nil,
				authToken: "",
				response: struct {
					Message string `json:"message"`
				}{
					Message: "Unauthorized",
				},
			},
		},
		{
			name: "login failed",
			args: args{
				setToken:  true,
				authToken: wantToken,
			},
			status: 401,
			want: want{
				user:      nil,
				authToken: wantToken,
				response: struct {
					Message string `json:"message"`
				}{
					Message: "Unauthorized",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authManagerMock := &AuthManagerMock{}
			authManagerMock.LoginSessionFunc = func(
				ctx context.Context, token string,
			) (*models.User, error) {
				if tt.status == 200 {
					if token != tt.args.authToken {
						t.Errorf("token is invalid. got: %s, want: %s", token, tt.args.authToken)
					}
					return tt.want.user, nil
				}
				return nil, errors.New("failed to login")
			}

			sut := NewAuthSessionLoginHandler(authManagerMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)
			r = testutil.SetLoggerContextToRequest(t, r)

			if tt.args.setToken {
				r.Header.Set("Authorization", "Bearer "+tt.args.authToken)
			}
			if tt.args.headerFormatInvalid {
				r.Header.Set("Authorization", tt.args.authToken)
			}

			sut.ServeHTTP(w, r)

			wb, err := json.Marshal(tt.want.response)
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

func Test_AuthLogoutHandler(t *testing.T) {
	t.Skip("TODO")

	type response struct {
		Message string `json:"message"`
	}
	wantToken := "authToken"
	type want struct {
		isSetCookie bool
		response
	}
	tests := []struct {
		name   string
		args   interface{}
		status int
		want   want
	}{
		{
			name:   "success",
			args:   nil,
			status: 200,
			want: want{
				isSetCookie: false,
				response: response{
					Message: "success",
				},
			},
		},
	}

	cookie := NewCookieManager("test", "https://test.example.com")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewAuthLogoutHandler(cookie)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)
			r = testutil.SetLoggerContextToRequest(t, r)

			cookie.SetCookie(w, "authToken", wantToken)
			fmt.Println(w.Header().Get("Set-Cookie"))

			sut.ServeHTTP(w, r)

			fmt.Println(w.Header().Get("Set-Cookie"))
			wb, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("cannot marshal want: %v", err)
			}

			resp := w.Result()
			if err := testutil.AssertResponse(t, resp, tt.status, wb); err != nil {
				t.Error(err)
			}

			parser := &http.Request{Header: http.Header{"Cookie": resp.Header["Set-Cookie"]}}
			if !tt.want.isSetCookie {
				token, err := parser.Cookie("authToken")
				fmt.Println("### token")
				fmt.Println(token)
				fmt.Println(token.Value)
				if err != nil {
					t.Errorf("failed to get cookie: %v", err)
				}
				if token.Value == wantToken {
					t.Errorf("cookie is not cleared")
				}
			}
		})
	}
}
