package app

import (
	"log"
	"net/http"
	"strconv"
	"tradingchat/pkg/chat"
	"tradingchat/pkg/model"

	"github.com/gorilla/mux"
)

func (app *App) handleEnterRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomStr, ok := params["room"]
	if !ok {
		http.Error(w, "missing room parameter", http.StatusBadRequest)
		return
	}

	// Parameter is always a number due to route configuration
	room, _ := strconv.Atoi(roomStr)

	if room < 0 || room >= len(app.Bgs) {
		http.Error(w, "room out of bounds", http.StatusBadRequest)
		return
	}

	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	user := r.Context().Value(UserContextKey("user")).(model.User)

	chat.HandleConnection(conn, app.Bgs[room], &user)
}
