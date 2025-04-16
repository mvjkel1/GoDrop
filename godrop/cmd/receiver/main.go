package main

import (
	"fmt"
	"net/http"

	"godrop/internal/transfer"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleReceive(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	err = transfer.ReceiveFile(conn)
	if err != nil {
		fmt.Println("Error receiving file:", err)
	} else {
		fmt.Println("File received successfully!")
	}
}

func main() {
	http.HandleFunc("/ws", handleReceive)
	fmt.Println("Receiver is listening on :8080")
	http.ListenAndServe(":8080", nil)
}
