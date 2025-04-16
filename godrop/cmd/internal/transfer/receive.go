package transfer

import (
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

func ReceiveFile(conn *websocket.Conn, passphrase string) error {
	_, name, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	cleanName := filepath.Base(string(name))
	f, err := os.Create("received_" + cleanName)
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		_, nonce, err := conn.ReadMessage()
		if err != nil {
			break
		}

		_, encryptedData, err := conn.ReadMessage()
		if err != nil {
			break
		}

		plaintext, err := Decrypt(nonce, encryptedData, passphrase)
		if err != nil {
			return err
		}

		_, err = f.Write(plaintext)
		if err != nil {
			return err
		}
	}

	return nil
}
