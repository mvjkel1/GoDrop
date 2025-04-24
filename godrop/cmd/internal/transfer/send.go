package transfer

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/schollz/progressbar/v3"
)

func SendFile(filePath, peerAddr, passphrase string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	url := "ws://" + peerAddr + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(stat.Name()))
	if err != nil {
		return fmt.Errorf("send filename: %w", err)
	}

	bar := progressbar.DefaultBytes(
		stat.Size(),
		"Sending",
	)

	buffer := make([]byte, 8) // []byte, 8 to see the progress bar actually loading, to be changed
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			nonce, encryptedData, err := Encrypt(buffer[:n], passphrase)
			if err != nil {
				return fmt.Errorf("encrypt: %w", err)
			}

			err = conn.WriteMessage(websocket.BinaryMessage, nonce)
			if err != nil {
				return fmt.Errorf("send nonce: %w", err)
			}
			err = conn.WriteMessage(websocket.BinaryMessage, encryptedData)
			if err != nil {
				return fmt.Errorf("send encrypted data: %w", err)
			}

			bar.Add(n)
		}
		if err != nil {
			break
		}
	}
	return nil
}
