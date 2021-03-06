package app

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"tradingchat/pkg/chat"

	"github.com/streadway/amqp"
)

func (app *App) ConsumeQueueMessages(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	return ch.Consume("stocks.q", "", true, false, false, false, nil)
}

func (app *App) HandleConsumedMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		fmt.Printf("got message: %s\n", msg.Body)
		var msgStruct struct {
			Message  string `json:"message"`
			Error    bool   `json:"hasErr"`
			CallerID string `json:"caller_id"`
		}
		err := json.Unmarshal(msg.Body, &msgStruct)
		if err != nil {
			continue
		}
		bgIndex, err := strconv.Atoi(msgStruct.CallerID)
		if err != nil || bgIndex >= len(app.Bgs) || bgIndex < 0 {
			fmt.Printf("invalid broadcast group: %s\n", msgStruct.CallerID)
			continue
		}

		msgChat := chat.Message{
			User:      "stock_bot",
			Role:      "bot",
			Content:   msgStruct.Message,
			Error:     msgStruct.Error,
			Timestamp: time.Now().Format("15:04:05"),
		}

		bytes, err := json.Marshal(msgChat)
		if err != nil {
			continue
		}

		app.Bgs[bgIndex].Get(bytes)

	}
}
