package discovery

import (
	"fmt"
	"net"
	"time"
)

func DiscoverPeers(broadcastPort string, timeout time.Duration) ([]string, error) {
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return nil, err
	}
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	broadcastAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:"+broadcastPort)
	if err != nil {
		return nil, err
	}

	_, err = conn.WriteToUDP([]byte("GODROP_DISCOVERY"), broadcastAddr)
	if err != nil {
		return nil, err
	}

	var peers []string
	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			break
		}
		response := string(buffer[:n])
		fmt.Printf("Response from %s: %s\n", addr.String(), response)
		peers = append(peers, response)
	}

	return peers, nil
}
