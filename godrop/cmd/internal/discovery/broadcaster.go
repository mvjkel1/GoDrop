package discovery

import (
	"fmt"
	"net"
	"time"
)

func DiscoverPeers(broadcastPort string, timeout time.Duration) ([]string, error) {
	conn, err := setupConnection(timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = sendDiscoveryBroadcast(conn, broadcastPort)
	if err != nil {
		return nil, err
	}

	peers := readResponses(conn)
	return peers, nil
}

func setupConnection(timeout time.Duration) (*net.UDPConn, error) {
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return nil, err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
		conn.Close()
		return nil, err
	}

	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func sendDiscoveryBroadcast(conn *net.UDPConn, broadcastPort string) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:"+broadcastPort)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP([]byte("GODROP_DISCOVERY"), broadcastAddr)
	return err
}

func readResponses(conn *net.UDPConn) []string {
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

	return peers
}
