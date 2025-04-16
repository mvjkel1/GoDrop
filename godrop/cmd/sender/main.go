package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"godrop/cmd/internal/discovery"
	"godrop/cmd/internal/transfer"
)

func main() {
	fmt.Println("Searching for peers on the network...")
	peers, err := discovery.DiscoverPeers("9999", 2*time.Second)
	if err != nil || len(peers) == 0 {
		fmt.Println("No peers found.")
		return
	}

	fmt.Println("Found peers:")
	for i, p := range peers {
		fmt.Printf("  [%d] %s\n", i+1, p)
	}

	fmt.Print("Choose a peer to send to [1]: ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)

	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil || choice < 1 || choice > len(peers) {
		choice = 1
	}
	peerAddr := peers[choice-1]

	fmt.Print("Enter passphrase: ")
	passphrase, _ := reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)
	fmt.Println("Passphrase accepted.")

	fmt.Print("Enter path to file: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File not found:", filePath)
		return
	}

	fmt.Println("Sending file to", peerAddr)
	err = transfer.SendFile(filePath, peerAddr, passphrase)
	if err != nil {
		fmt.Println("Failed to send file:", err)
	} else {
		fmt.Println("File sent successfully!")
	}
}
