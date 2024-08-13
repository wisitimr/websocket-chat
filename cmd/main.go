package main

import (
	"context"
	"log"
	server "websocket-chat/internal"
	"websocket-chat/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("parse .env error")
	}
	server, err := server.NewHttpServer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(ctx)
}
