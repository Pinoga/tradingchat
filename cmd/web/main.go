package main

import (
	"log"
	"tradingchat/pkg/app"
)

func main() {
	config := app.AppConfig{
		LenBgs:          1,
		CookieSecretKey: "abcdef",
		Port:            ":8080",
	}
	application := app.NewApp()
	application.Initialize(config)
	err := application.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
