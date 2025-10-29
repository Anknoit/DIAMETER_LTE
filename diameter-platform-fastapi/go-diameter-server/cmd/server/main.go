package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"diam/internal/api"
)

func main() {
	log.Println("diameter server scaffold starting")

	// TODO: replace with real config loader
	go func() {
		mux := http.NewServeMux()
		api.RegisterHandlers(mux)
		log.Println("starting management API on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("management server error: %v", err)
		}
	}()

	// Placeholder: real Diameter listener (TCP/TLS) to be implemented in internal/diameter
	log.Println("Diameter protocol engine not implemented in scaffold (see internal/diameter)")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down diameter server scaffold")
	// graceful shutdown logic would go here
}
