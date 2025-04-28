package infrastructure

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisKVS struct {
	cli           *redis.Client
	expirationSec int
}

type RedisOption func(*redis.Options)

func NewRedisKVS(
	ctx context.Context,
	host string,
	port int64,
	username string,
	password string,
	expirationSec int,
	enableTLS bool,
) (*RedisKVS, error) {
	redisOpt := &redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: username,
		Password: password,
	}

	if enableTLS {
		redisOpt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	cli := redis.NewClient(redisOpt)
	if res := cli.Ping(ctx); res.Err() != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", res.Err())
	}
	kvs := &RedisKVS{
		cli:           cli,
		expirationSec: expirationSec,
	}
	return kvs, nil
}

func (r *RedisKVS) Save(ctx context.Context, key string, value string) error {
	ret := r.cli.Set(ctx, key, value, time.Second*time.Duration(r.expirationSec))
	if ret.Err() != nil {
		return fmt.Errorf("failed to set key: %w", ret.Err())
	}
	return nil
}

func (r *RedisKVS) Load(ctx context.Context, key string) (string, error) {
	ret := r.cli.Get(ctx, key)
	if ret.Err() != nil {
		return "", fmt.Errorf("failed to get key: %w", ret.Err())
	}
	return ret.Val(), nil
}
