package main

import (
	"log"

	"github.com/grandcat/zeroconf"
)

type BonjourHandle struct {
	server *zeroconf.Server
}

func startBonjour(cfg AppConfig) (*BonjourHandle, error) {
	txt := []string{
        "app=vms-api",                // nama aplikasi
        "ver=1.0",                     // versi API
        "path=/api",                   // base path HTTP
        "desc=VMS HTTP API Server",    // deskripsi singkat
    }
	srv, err := zeroconf.Register(
		cfg.ServiceName,
		cfg.ServiceType,
		cfg.ServiceDomain,
		cfg.Port,
		txt,
		nil,
	)
	if err != nil {
		return nil, err
	}
	log.Printf("Bonjour advertised: %s.%s%s on port %d",
		cfg.ServiceName, cfg.ServiceType, cfg.ServiceDomain, cfg.Port)
	return &BonjourHandle{server: srv}, nil
}

func (b *BonjourHandle) Stop() {
	if b.server != nil {
		b.server.Shutdown()
	}
}
