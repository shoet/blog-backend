package jwt_service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type KVSer interface {
	Save(ctx context.Context, key string, value string) error
	Load(ctx context.Context, key string) (string, error)
}

type JWTService struct {
	kvs               KVSer
	clocker           clocker.Clocker
	secretKey         []byte
	tokenExpiresInSec int
}

func NewJWTService(
	kvs KVSer,
	clocker clocker.Clocker,
	secretKey []byte,
	tokenExpiresInSec int,
) *JWTService {
	return &JWTService{
		kvs:               kvs,
		clocker:           clocker,
		secretKey:         secretKey,
		tokenExpiresInSec: tokenExpiresInSec,
	}
}

func (j *JWTService) GenerateToken(ctx context.Context, u *models.User) (string, error) {
	uuid := uuid.New().String()
	claims := jwt.RegisteredClaims{
		ID:       uuid,
		Subject:  "blog",
		IssuedAt: jwt.NewNumericDate(j.clocker.Now()),
		ExpiresAt: jwt.NewNumericDate(
			j.clocker.Now().Add(time.Duration(j.tokenExpiresInSec) * time.Second)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	userIdStr := strconv.Itoa(int(u.Id))
	if err := j.kvs.Save(context.Background(), uuid, userIdStr); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}
	return ss, nil
}

var ErrSessionNotFound = errors.New("session is not found")

func (j *JWTService) VerifyToken(ctx context.Context, token string) (models.UserId, error) {
	// parse token
	parsed, err := jwt.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	claims := parsed.Claims.(*jwt.RegisteredClaims)

	// check session kvs
	v, err := j.kvs.Load(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrSessionNotFound
		}
		return 0, fmt.Errorf("failed to load token: %w", err)
	}
	userId, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("failed to convert user id: %w", err)
	}
	return models.UserId(userId), nil
}
