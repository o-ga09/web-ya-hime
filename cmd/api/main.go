package main

import (
	"context"
	"log"

	"github.com/o-ga09/web-ya-hime/internal/server"
	"github.com/o-ga09/web-ya-hime/pkg/config"
)

func main() {
	ctx, err := config.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(ctx)
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
