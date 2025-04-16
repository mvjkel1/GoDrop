package transfer

import (
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func SendFile(filePath, peerAddr string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	url := "ws://" + peerAddr + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(file.Name()))
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
