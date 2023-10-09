package main

import (
	"context"
	"log"

	"github.com/shoet/blog/handlers"
)

func main() {
	server, err := handlers.NewServer(3000)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
