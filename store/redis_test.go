package store

import (
	"context"
	"testing"
)

func Test_NewRedisKVS(t *testing.T) {
	ctx := context.Background()
	_, err := NewRedisKVS(ctx, "127.0.0.1", 6379, "default", "redispw", 10, false)
	if err != nil {
		t.Fatalf("failed to create redis kvs: %v", err)
	}
}

func Test_Save(t *testing.T) {
	ctx := context.Background()
	kvs, err := NewRedisKVS(ctx, "127.0.0.1", 6379, "default", "redispw", 10, false)
	if err != nil {
		t.Fatalf("failed to create redis kvs: %v", err)
	}

	if err := kvs.Save(ctx, "test", "test"); err != nil {
		t.Fatalf("failed to save: %v", err)
	}
}

func Test_Load(t *testing.T) {
	ctx := context.Background()
	kvs, err := NewRedisKVS(ctx, "127.0.0.1", 6379, "default", "redispw", 300, false)
	if err != nil {
		t.Fatalf("failed to create redis kvs: %v", err)
	}

	want := "test want"

	if err := kvs.Save(ctx, "test", want); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	ret, err := kvs.Load(ctx, "test")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if ret != want {
		t.Errorf("want %s, got %s", want, ret)
	}
}
