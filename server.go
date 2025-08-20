package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type Message struct {
	ID         int64  `json:"id"`
	RemoteAddr string `json:"remote_addr"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}


func parseContent(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return ""
	}
	var msg Message
	if json.Unmarshal([]byte(line), &msg) == nil && strings.TrimSpace(msg.Content) != "" {
		return msg.Content
	}
	return trimmed
}

func startTCPServer(cfg AppConfig, db *AppDB) (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("accept error: %v", err)
				return
			}
			go handleConn(conn, db)
		}
	}()

	log.Printf("TCP server listening on :%d", cfg.Port)
	return ln, nil
}

func handleConn(conn net.Conn, db *AppDB) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	reader := bufio.NewScanner(conn)
	buf := make([]byte, 0, 64*1024)
	reader.Buffer(buf, 1024*1024)

	for reader.Scan() {
		line := reader.Text()
		content := parseContent(line)
		if content != "" {
			if err := db.InsertMessage(remote, content); err != nil {
				log.Printf("DB insert error: %v", err)
				fmt.Fprintf(conn, "ERR: %v\n", err)
				continue
			}
			fmt.Fprintln(conn, "OK")
		}
	}
	if err := reader.Err(); err != nil {
		log.Printf("read error from %s: %v", remote, err)
	}
}
