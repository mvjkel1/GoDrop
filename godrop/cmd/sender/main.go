package main

import (
	"fmt"
	"godrop/internal/transfer"
)

func main() {
	filePath := "example.txt"
	peerAddr := "192.168.1.100:8080"

	err := transfer.SendFile(filePath, peerAddr)
	if err != nil {
		fmt.Println("Failed to send file:", err)
	} else {
		fmt.Println("File sent successfully!")
	}
}
