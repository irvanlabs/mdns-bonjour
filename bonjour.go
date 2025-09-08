package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

type BonjourHandle struct {
	server   *zeroconf.Server
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// baca default route dari /proc/net/route
func getDefaultIfaceName() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		// kolom 2 = Destination (00000000 berarti default route)
		if len(fields) >= 2 && fields[1] == "00000000" {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("no default route found")
}

func getIfaceAndIP() (*net.Interface, net.IP, error) {
	ifaceName, err := getDefaultIfaceName()
	if err != nil {
		return nil, nil, err
	}

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				return iface, ip4, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("no IPv4 found on %s", ifaceName)
}

func startBonjour(cfg AppConfig) *BonjourHandle {
	handle := &BonjourHandle{stopChan: make(chan struct{})}

	handle.wg.Add(1)
	go func() {
		defer handle.wg.Done()
		var lastIP string

		for {
			select {
			case <-handle.stopChan:
				if handle.server != nil {
					handle.server.Shutdown()
					handle.server = nil
				}
				return
			default:
			}

			iface, ip, err := getIfaceAndIP()
			if err != nil {
				log.Println("Network belum siap:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if ip.String() != lastIP {
				// IP berubah → re-register Zeroconf
				if handle.server != nil {
					log.Printf("🔄 Network berubah (%s → %s), restart announce...", lastIP, ip.String())
					handle.server.Shutdown()
					handle.server = nil
				}

				txt := []string{
					"app=vms-api",
					"ver=1.0",
					"path=/api",
					"desc=VMS HTTP API Server",
				}

				srv, err := zeroconf.Register(
					cfg.ServiceName,
					cfg.ServiceType,
					cfg.ServiceDomain,
					cfg.Port,
					txt,
					[]net.Interface{*iface},
				)
				if err != nil {
					log.Println("❌ Register gagal:", err)
					time.Sleep(5 * time.Second)
					continue
				}

				handle.server = srv
				lastIP = ip.String()
				log.Printf("✅ Berhasil announce di %s (%s)", iface.Name, ip.String())
			}

			time.Sleep(10 * time.Second)
		}
	}()

	return handle
}

func (b *BonjourHandle) Stop() {
	close(b.stopChan)
	b.wg.Wait()
}
