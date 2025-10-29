package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: load config, init logger, init stores, start diameter listener and http api
	log.Println("diameter server starting (scaffold)")

	// placeholder graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down")
}
