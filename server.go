package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	ID         int64  `json:"id"`
	RemoteAddr string `json:"remote_addr"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}


func startTCPServer(cfg AppConfig) (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	log.Printf("TCP server listening on :%d", cfg.Port)
	return ln, nil
}
