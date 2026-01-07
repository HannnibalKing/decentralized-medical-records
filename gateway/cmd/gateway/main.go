package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/example/med/gateway/internal/config"
	"github.com/example/med/gateway/internal/server"
)

func main() {
	cfg := config.FromEnv()

	srv := server.New(cfg)
	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("gateway listening on %s", cfg.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
	_ = os.Stdout
}
