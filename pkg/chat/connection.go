package chat

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"tradingchat/internal/model"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleConnection(conn *websocket.Conn, bg *BroadcastGroup) {
	client := Client{
		IncomingMessages: make(chan []byte),
	}

	defer func() {
		bg.Leave(client)
		conn.Close()
	}()

	bg.Enter(client)

	go handleClientMessage(conn, bg)
	go handleBroadcastMessage(conn, &client)
}

func handleClientMessage(conn *websocket.Conn, bg *BroadcastGroup) {
	for {
		_, r, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			break
		}

		bg.Get(r)
	}
}

func handleBroadcastMessage(conn *websocket.Conn, client *Client) {
	for msgContent := range client.IncomingMessages {
		msgStruct := model.Message{
			Content: string(msgContent),
		}
		payload, err := json.Marshal(&msgStruct)
		if err != nil {
			return
		}
		write(payload, conn)
	}
}

func write(msg []byte, conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.TextMessage, msg)
}
