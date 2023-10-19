package interfaces

import (
	"context"
)

type KVSer interface {
	Save(ctx context.Context, key string, value string) error
	Load(ctx context.Context, key string) (string, error)
}
