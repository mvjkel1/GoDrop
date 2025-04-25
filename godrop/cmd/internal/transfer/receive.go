package transfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

func ReceiveFile(conn *websocket.Conn, passphrase string) error {
	file, err := prepareOutputFile(conn)
	if err != nil {
		return err
	}
	defer file.Close()

	return receiveEncryptedStream(conn, file, passphrase)
}

func prepareOutputFile(conn *websocket.Conn) (*os.File, error) {
	_, nameMsg, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read filename: %w", err)
	}

	cleanName := filepath.Base(string(nameMsg))
	outPath := "received_" + cleanName

	file, err := os.Create(outPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}

	return file, nil
}

func receiveEncryptedStream(conn *websocket.Conn, file *os.File, passphrase string) error {
	for {
		_, nonce, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) || err == io.EOF {
				break
			}
			return fmt.Errorf("read nonce: %w", err)
		}

		_, encryptedData, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read encrypted data: %w", err)
		}

		plaintext, err := Decrypt(nonce, encryptedData, passphrase)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}

		if _, err := file.Write(plaintext); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
	}
	return nil
}
