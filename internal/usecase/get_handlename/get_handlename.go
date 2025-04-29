package get_handlename

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type KVS interface {
	Load(ctx context.Context, key string) (*string, error)
	Save(ctx context.Context, key string, value string) error
}

type Usecase struct {
	kvs KVS
}

func NewUsecase(kvs KVS) *Usecase {
	return &Usecase{kvs: kvs}
}

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId, ip string) (string, error) {
	saltKey := fmt.Sprintf("handlename.salt.%d", blogId)

	salt, err := u.kvs.Load(ctx, saltKey)
	if err != nil {
		return "", fmt.Errorf("failed to get salt: %w", err)
	}
	if salt == nil {
		newSalt := uuid.NewString()
		if err := u.kvs.Save(ctx, saltKey, newSalt); err != nil {
			return "", fmt.Errorf("failed to save salt: %w", err)
		}
		salt = &newSalt
	}
	source := fmt.Sprintf("%d.%s", blogId, *salt)
	h := sha256.New()
	h.Write([]byte(source))
	hash := fmt.Sprintf("%x", h.Sum(nil))[:10]
	handleName := strings.ToUpper(hash)
	return handleName, nil
}
