package main

import (
	"context"
	"log"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/interfaces"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}
	server, err := interfaces.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := server.Run(ctx); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
