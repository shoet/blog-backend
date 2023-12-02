package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/services"
)

func Test_JWTService_GenerateToken(t *testing.T) {
	type args struct {
		user *models.User
	}

	type want struct {
		user *models.User
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name:    "success",
			args:    args{user: &models.User{Id: 1}},
			want:    want{user: &models.User{Id: 1}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			kvsMock := &services.KVSerMock{}
			kvsMock.SaveFunc = func(ctx context.Context, key string, value string) error {
				if value != strconv.Itoa(int(tt.args.user.Id)) {
					t.Fatalf("failed want user id: %v, got: %v", tt.want.user.Id, value)
				}
				return nil
			}
			clockerMock := &clocker.FiexedClocker{}

			testSecret := "12345678"
			testTokenExpiresInSec := 60
			sut := NewJWTer(kvsMock, clockerMock, []byte(testSecret), testTokenExpiresInSec)

			token, err := sut.GenerateToken(ctx, tt.args.user)
			if err != nil {
				t.Fatalf("failed generate token: %v", err)
			}
			_ = token
		})
	}
}

func Test_JWTService_VerifyToken(t *testing.T) {
	type args struct {
		user             *models.User
		tokenExpireInSec int
	}

	type want struct {
		user *models.User
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name: "success",
			args: args{
				user:             &models.User{Id: 1},
				tokenExpireInSec: 60,
			},
			want:    want{user: &models.User{Id: 1}},
			wantErr: nil,
		},
		{
			name: "failed expired token",
			args: args{
				user:             &models.User{Id: 1},
				tokenExpireInSec: 0,
			},
			want:    want{user: &models.User{Id: 1}},
			wantErr: fmt.Errorf("token is expired"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			kvsMock := &services.KVSerMock{}
			var uuid string
			kvsMock.SaveFunc = func(ctx context.Context, key string, value string) error {
				uuid = key // 発行されたuuidを保持しておく
				if value != strconv.Itoa(int(tt.args.user.Id)) {
					t.Fatalf("failed want user id: %v, got: %v", tt.want.user.Id, value)
				}
				return nil
			}
			kvsMock.LoadFunc = func(ctx context.Context, key string) (string, error) {
				if key != uuid {
					t.Fatalf("failed want uuid: %v, got: %v", uuid, key)
				}
				return strconv.Itoa(int(tt.args.user.Id)), nil
			}
			clockerMock := &clocker.RealClocker{}
			testSecret := "12345678"
			testTokenExpiresInSec := tt.args.tokenExpireInSec

			sut := NewJWTer(kvsMock, clockerMock, []byte(testSecret), testTokenExpiresInSec)

			token, err := sut.GenerateToken(ctx, tt.args.user)
			if err != nil {
				t.Fatalf("failed generate token: %v", err)
			}

			userId, err := sut.VerifyToken(ctx, token)
			if err != nil {
				if strings.Contains(err.Error(), tt.wantErr.Error()) {
					return
				} else {
					t.Fatalf("failed verify token: %v", err)
					return
				}
			}

			if userId != tt.want.user.Id {
				t.Fatalf("failed want user id: %v, got: %v", tt.want.user.Id, userId)
			}
		})
	}

}
