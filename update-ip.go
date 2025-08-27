package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type DevicePayload struct {
	ServerId     string `json:"server_id"`
	LocalIp       string `json:"local_ip"`
	PublicIp     string `json:"public_ip"`
	Hostname string `json:"hostname"`
	Time     string `json:"time"`
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}

func getPublicIP() string {
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "error"
	}
	return strings.TrimSpace(string(body)) // hapus newline
}



func announce(cfg AppConfig)  {
	payload := deviceIp(cfg)
	resp, err := http.Post("https://webhook.site/6b0b379d-0cd2-4029-a8a8-b89d66eb018c", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("❌ Gagal kirim:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("✅ Berhasil announce:", string(payload), "status:", resp.Status)
}

func deviceIp(cfg AppConfig) []byte {
	payload := DevicePayload{
		ServerId:    cfg.ServerId,
		LocalIp:     getLocalIP(),
		PublicIp:    getPublicIP(),
		Hostname: getHostname(),
		Time:     time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)

	return body
}

// Fungsi ambil hostname
func getHostname() string {
	name, err := net.LookupAddr(getLocalIP())
	if err == nil && len(name) > 0 {
		return name[0]
	}
	// fallback
	hostname, _ := net.LookupHost("localhost")
	if len(hostname) > 0 {
		return hostname[0]
	}
	return "unknown"
}



func updateIp(cfg AppConfig) {
    c := cron.New(cron.WithSeconds())
    c.AddFunc("0 0 * * * *", func() { announce(cfg) })
	c.Start()
	announce(cfg)
}
