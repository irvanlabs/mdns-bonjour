package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := LoadConfig()

	db := NewAppDB(cfg.DBPath)
	defer db.Close()

	ln, err := startTCPServer(cfg, db)
	if err != nil {
		log.Fatalf("start tcp server: %v", err)
	}
	defer ln.Close()

	bonjour, err := startBonjour(cfg)
	if err != nil {
		log.Fatalf("start bonjour: %v", err)
	}
	defer bonjour.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Service running... Press Ctrl+C to stop.")
	<-ctx.Done()
	log.Println("Shutting down...")
}
