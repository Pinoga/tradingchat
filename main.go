package main

import (
	"net/http"
	"tradingchat/chat/broadcast"
	"tradingchat/chat/connection"
)

func main() {
	address := ":8080"
	group := broadcast.NewBroadCastGroup()
	go group.HandleBroadcasts()
	http.HandleFunc("/chat", func(rw http.ResponseWriter, r *http.Request) {
		connection.HandleConnection(rw, r, group)
	})
	http.ListenAndServe(&address, nil)
}
