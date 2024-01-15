package jwt_service_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/services/jwt_service"
	"github.com/stretchr/testify/mock"
)

type KVSerMock struct {
	mock.Mock
}

func (m *KVSerMock) Save(ctx context.Context, key string, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *KVSerMock) Load(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

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
			kvsMock := &KVSerMock{}
			userIdStr := strconv.Itoa(int(tt.want.user.Id))
			kvsMock.On("Save", mock.Anything, mock.AnythingOfType("string"), userIdStr).Return(nil)
			clockerMock := &clocker.FiexedClocker{}

			testSecret := "12345678"
			testTokenExpiresInSec := 60
			sut := jwt_service.NewJWTService(kvsMock, clockerMock, []byte(testSecret), testTokenExpiresInSec)

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
			kvsMock := &KVSerMock{}
			userIdStr := strconv.Itoa(int(tt.want.user.Id))
			kvsMock.On("Save", mock.Anything, mock.AnythingOfType("string"), userIdStr).Return(nil)
			kvsMock.On("Load", mock.Anything, mock.AnythingOfType("string")).Return(userIdStr, nil)

			clockerMock := &clocker.RealClocker{}
			testSecret := "12345678"
			testTokenExpiresInSec := tt.args.tokenExpireInSec

			sut := jwt_service.NewJWTService(kvsMock, clockerMock, []byte(testSecret), testTokenExpiresInSec)

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
