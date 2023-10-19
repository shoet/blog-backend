package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/models"
)

type JWTService struct {
	kvs     KVSer
	cfg     *config.Config
	clocker clocker.Clocker
}

func NewJWTService(kvs KVSer, cfg *config.Config, clocker clocker.Clocker) *JWTService {
	return &JWTService{
		kvs:     kvs,
		cfg:     cfg,
		clocker: clocker,
	}
}

func (j *JWTService) GenerateToken(ctx context.Context, u *models.User) (string, error) {
	uuid := uuid.New().String()
	claims := jwt.RegisteredClaims{
		ID:        uuid,
		Subject:   "blog",
		IssuedAt:  jwt.NewNumericDate(j.clocker.Now()),
		ExpiresAt: jwt.NewNumericDate(j.clocker.Now().Add(time.Duration(j.cfg.JWTExpiresInSec) * time.Second)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(j.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	userIdStr := strconv.Itoa(int(u.Id))
	if err := j.kvs.Save(context.Background(), uuid, userIdStr); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}
	return ss, nil
}

func (j *JWTService) VerifyToken(ctx context.Context, token string) (models.UserId, error) {
	// parse token
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims := parsed.Claims.(*jwt.RegisteredClaims)

	// kvs
	v, err := j.kvs.Load(ctx, claims.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to load token: %w", err)
	}
	userId, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("failed to convert user id: %w", err)
	}
	return models.UserId(userId), nil
}
