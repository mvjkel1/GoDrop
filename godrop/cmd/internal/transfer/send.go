package transfer

import (
	"os"

	"github.com/gorilla/websocket"
)

func SendFile(filePath, peerAddr, passphrase string) error {
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

		nonce, encryptedData, err := Encrypt(buffer[:n], passphrase)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.BinaryMessage, nonce)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.BinaryMessage, encryptedData)
		if err != nil {
			return err
		}
	}

	return nil
}
