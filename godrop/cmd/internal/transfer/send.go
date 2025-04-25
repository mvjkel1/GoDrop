package transfer

import (
	"fmt"
	"io"
	"os"

	"github.com/gorilla/websocket"
	"github.com/schollz/progressbar/v3"
)

func SendFile(filePath, peerAddr, passphrase string) error {
	file, stat, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	conn, err := connectWebSocket(peerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendFileName(conn, stat.Name()); err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(stat.Size(), "Sending")
	return streamFile(conn, file, passphrase, bar)
}

func openFile(path string) (*os.File, os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("stat file: %w", err)
	}
	return file, stat, nil
}

func connectWebSocket(addr string) (*websocket.Conn, error) {
	url := "ws://" + addr + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return conn, nil
}

func sendFileName(conn *websocket.Conn, fileName string) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte(fileName))
	if err != nil {
		return fmt.Errorf("send filename: %w", err)
	}
	return nil
}

func streamFile(conn *websocket.Conn, file *os.File, passphrase string, bar *progressbar.ProgressBar) error {
	buffer := make([]byte, 8) // small for testing, increase later

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			nonce, encryptedData, err := Encrypt(buffer[:n], passphrase)
			if err != nil {
				return fmt.Errorf("encrypt: %w", err)
			}

			if err := conn.WriteMessage(websocket.BinaryMessage, nonce); err != nil {
				return fmt.Errorf("send nonce: %w", err)
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, encryptedData); err != nil {
				return fmt.Errorf("send encrypted data: %w", err)
			}
			bar.Add(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
	}
	return nil
}
