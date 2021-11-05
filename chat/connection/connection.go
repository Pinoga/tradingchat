package connection

import (
	"io"
	"log"
	"net/http"
	"tradingchat/chat"
	"tradingchat/chat/broadcast"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func HandleConnection(w http.ResponseWriter, r *http.Request, bg broadcast.BroadcastGroup) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	client := broadcast.Client{
		IncomingMessages: make(chan<- chat.Message),
	}
	bg.Enter(client)

	for {
		messageType, r, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			break
		}

		msg := chat.Message{
			Content: r.Read(),
		}

		bg.Get()
	}

	bg.Leave(client)
}

func write(msg chat.Message, conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
}
