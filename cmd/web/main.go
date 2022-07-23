package main

import (
	"log"
	"os"
	"strconv"
	"tradingchat/pkg/app"
)

func main() {
	lenBgs, err := strconv.Atoi(os.Getenv("BROADCAST_GROUPS"))
	if err != nil {
		log.Fatalf("couldn't parse env variable. %v", err)
	}
	config := app.AppConfig{
		LenBgs:          lenBgs,
		CookieSecretKey: os.Getenv("COOKIE_SECRET_KEY"),
		Port:            os.Getenv("APP_PORT"),
		DatabaseURI:     os.Getenv("DB_URI"),
		DatabaseName:    os.Getenv("DB_NAME"),
		RabbitMQURI:     os.Getenv("RABBITMQ_URI"),
		RedisHost:       os.Getenv("REDIS_HOST"),
	}
	application := app.NewApp()
	defer application.Stop()
	application.Initialize(config)
}
