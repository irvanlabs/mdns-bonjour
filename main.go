package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := LoadConfig()

	ln, err := startTCPServer(cfg)
	if err != nil {
		log.Fatalf("start tcp server: %v", err)
	}
	defer ln.Close()

	bonjour, err := startBonjour(cfg)
	if err != nil {
		log.Fatalf("start bonjour: %v", err)
	}
	defer bonjour.Stop()

	router := Router()
	go func() {
		log.Printf("HTTP server listening on :%d", cfg.ApiPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ApiPort), router); err != nil {
			log.Fatalf("start http server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Service running... Press Ctrl+C to stop.")
	updateIp(cfg)
	<-ctx.Done()
	log.Println("Shutting down...")
}

