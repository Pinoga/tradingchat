package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func PublishToQueue(ch *amqp.Channel) {
	for resp := range Responses {
		fmt.Println("res")
		ch.Publish("", "stocks.q", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        resp,
		})
	}
}
