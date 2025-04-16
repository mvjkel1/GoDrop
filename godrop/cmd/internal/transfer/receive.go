package transfer

import (
	"io"
	"os"

	"github.com/gorilla/websocket"
)

func ReceiveFile(conn *websocket.Conn) error {
	_, name, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	file, err := os.Create("received_" + string(name))
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF || websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			return err
		}
		if msgType == websocket.BinaryMessage {
			_, err := file.Write(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
