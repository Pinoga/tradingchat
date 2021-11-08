package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tradingchat/pkg/chat"

	"github.com/gorilla/mux"
)

func (app *App) handleEnterRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entered upgrade")
	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	params := mux.Vars(r)

	roomStr, ok := params["room"]
	if !ok {
		http.Error(w, "missing room parameter", http.StatusBadRequest)
	}

	// Parameter is always a number due to route configuration
	room, _ := strconv.Atoi(roomStr)

	if room < 0 || room >= len(app.Bgs) {
		http.Error(w, "room out of bounds", http.StatusBadRequest)
	}

	chat.HandleConnection(conn, &app.Bgs[room])
}
