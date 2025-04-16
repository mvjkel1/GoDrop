package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"godrop/cmd/internal/discovery"
	"godrop/cmd/internal/transfer"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func makeReceiveHandler(passphrase string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		err = transfer.ReceiveFile(conn, passphrase)
		if err != nil {
			fmt.Println("Error receiving file:", err)
		} else {
			fmt.Println("File received successfully!")
		}
	}
}

func main() {
	fmt.Print("Enter passphrase: ")
	reader := bufio.NewReader(os.Stdin)
	passphrase, _ := reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)

	go func() {
		err := discovery.StartDiscoveryResponder("9999", "8080")
		if err != nil {
			fmt.Println("Discovery responder error:", err)
		}
	}()

	http.HandleFunc("/ws", makeReceiveHandler(passphrase))
	fmt.Println("Receiver is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
