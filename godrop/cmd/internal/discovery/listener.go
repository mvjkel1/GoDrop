package discovery

import (
	"fmt"
	"net"
	"strings"
)

func StartDiscoveryResponder(listenPort string, wsPort string) error {
	addr, err := net.ResolveUDPAddr("udp4", ":"+listenPort)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("UDP read error:", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		if message == "GODROP_DISCOVERY" {
			localIP := getLocalIP()
			response := fmt.Sprintf("%s:%s", localIP, wsPort)
			conn.WriteToUDP([]byte(response), remoteAddr)
			fmt.Printf("Discovery request received, replied to %s\n", remoteAddr)
		}
	}
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			if ip.IsLinkLocalUnicast() {
				continue
			}

			return ip.String()
		}
	}
	return ""
}
