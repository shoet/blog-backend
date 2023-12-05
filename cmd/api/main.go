package main

import (
	"context"
	"log"

	"github.com/shoet/blog/config"
	"github.com/shoet/blog/handlers"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}
	server, err := handlers.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := server.Run(ctx); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
