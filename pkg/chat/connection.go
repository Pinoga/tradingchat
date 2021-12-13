package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"tradingchat/pkg/service"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleConnection(conn *websocket.Conn, bg *BroadcastGroup, user *service.User) {
	client := Client{
		IncomingMessages: make(chan []byte),
		User:             *user,
	}

	bg.Enter(&client)
	fmt.Println(bg.ID, bg.clients)

	go handleClientMessage(conn, bg, &client)
	go handleBroadcastMessage(conn, &client)
}

func handleClientMessage(conn *websocket.Conn, bg *BroadcastGroup, client *Client) {
	defer func() {
		bg.Leave(client)
		conn.Close()
	}()

	for {
		_, r, err := conn.ReadMessage()
		fmt.Printf("read message:%s\n", string(r))
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			break
		}

		content := string(r)

		// TODO: Decouple this and create a message controller for commands
		if strings.HasPrefix(content, "/stock=") {
			err := SendCommandToBot(strings.TrimPrefix(content, "/stock="), bg.ID)
			if err != nil {
				respMsg := Message{
					User:      "system",
					Role:      "system",
					Content:   err.Error(),
					Timestamp: "",
				}
				byteMsg, err := json.Marshal(respMsg)
				if err != nil {
					fmt.Printf("error marshalling message to be sent: %v", err)
					return
				}
				write(byteMsg, conn)
			}

		}

		msg := Message{
			User:      client.User.Username,
			Role:      client.User.Role,
			Content:   content,
			Timestamp: time.Now().Format("15:04:05"),
		}
		byteMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("error marshalling message to be sent: %v", err)
			return
		}

		bg.Get(byteMsg)
	}
}

func handleBroadcastMessage(conn *websocket.Conn, client *Client) {
	for byteContent := range client.IncomingMessages {
		write(byteContent, conn)
	}
}

func write(msg []byte, conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.TextMessage, msg)
}
