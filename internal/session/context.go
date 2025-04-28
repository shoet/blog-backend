package session

import (
	"context"
	"errors"

	"github.com/shoet/blog/internal/infrastructure/models"
)

var UserIdContextKey = struct{}{}
var ErrUserIdNotFound = errors.New("user id is not found")

func SetUserId(ctx context.Context, userId models.UserId) context.Context {
	return context.WithValue(ctx, UserIdContextKey, userId)
}

func GetUserId(ctx context.Context) (models.UserId, error) {
	userId, ok := ctx.Value(UserIdContextKey).(models.UserId)
	if !ok {
		return 0, ErrUserIdNotFound
	}
	return userId, nil
}
