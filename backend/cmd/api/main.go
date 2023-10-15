package main

import (
	"context"
	"log"

	"github.com/shoet/blog/handlers"
)

func main() {
	ctx := context.Background()
	server, err := handlers.NewServer(ctx, 3000)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
