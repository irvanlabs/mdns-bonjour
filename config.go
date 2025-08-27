package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServiceName   string
	ServiceType   string
	ServiceDomain string
	Port          int
	DBPath        string
	ApiPort       int
	ServerId      string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func mustAtoi(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

func LoadConfig() AppConfig {
	// Load file .env kalau ada
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("warning: error loading .env: %v", err)
		}
	}

	return AppConfig{
		ServiceName:   getenv("SERVICE_NAME", "mdns-socket-svc"),
		ServiceType:   getenv("SERVICE_TYPE", "_http._tcp"),
		ServiceDomain: getenv("SERVICE_DOMAIN", "local."),
		Port:          mustAtoi(getenv("PORT", "49221"), 49221),
		DBPath:        getenv("DB_PATH", "data.db"),
		ApiPort:       mustAtoi(getenv("API_PORT", "8080"), 8080),
		ServerId:      getenv("SERVER_ID", "device-123-xyz"),
	}
}
