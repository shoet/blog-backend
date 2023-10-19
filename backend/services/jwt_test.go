package services

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/models"
)

func Test_JWTService_GenerateToken(t *testing.T) {
	ctx := context.Background()
	wantUserId := 1

	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed new config: %v", err)
	}
	kvs := KVSerMock{}
	kvs.SaveFunc = func(ctx context.Context, key string, value string) error {
		if value != strconv.Itoa(wantUserId) {
			t.Fatalf("failed want user id: %v, got: %v", wantUserId, value)
		}
		return nil
	}
	c := clocker.FiexedClocker{}
	sut := JWTService{
		kvs:     &kvs,
		cfg:     cfg,
		clocker: &c,
	}

	user := &models.User{
		Id: 1,
	}
	token, err := sut.GenerateToken(ctx, user)
	if err != nil {
		t.Fatalf("failed generate token: %v", err)
	}
	fmt.Println(token)
}

func Test_JWTService_VerifyToken(t *testing.T) {
	ctx := context.Background()
	wantUserId := 1

	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed new config: %v", err)
	}
	kvs := KVSerMock{}
	var uuid string
	kvs.SaveFunc = func(ctx context.Context, key string, value string) error {
		uuid = key
		if value != strconv.Itoa(wantUserId) {
			t.Fatalf("failed want user id: %v, got: %v", wantUserId, value)
		}
		return nil
	}
	kvs.LoadFunc = func(ctx context.Context, key string) (string, error) {
		if key != uuid {
			t.Fatalf("failed want uuid: %v, got: %v", uuid, key)
		}
		return strconv.Itoa(wantUserId), nil
	}
	c := clocker.RealClocker{}
	sut := JWTService{
		kvs:     &kvs,
		cfg:     cfg,
		clocker: &c,
	}

	user := &models.User{
		Id: 1,
	}
	token, err := sut.GenerateToken(ctx, user)
	if err != nil {
		t.Fatalf("failed generate token: %v", err)
	}

	userId, err := sut.VerifyToken(ctx, token)
	if err != nil {
		t.Fatalf("failed verify token: %v", err)
	}

	if userId != models.UserId(wantUserId) {
		t.Fatalf("failed want user id: %v, got: %v", wantUserId, userId)
	}

}
