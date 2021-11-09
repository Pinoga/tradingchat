package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/streadway/amqp"
)

type StockRequest struct {
	ID string `json:"ID"`
}

type RequestProcess struct {
	InternalID string
	CallerID   string
	StockCode  string
}

type StockResponse struct {
	Message  string `json:"message"`
	Error    bool   `json:"error"`
	CallerID string `json:"caller_id"`
}

var (
	Requests  = make(chan RequestProcess, 16)
	Responses = make(chan []byte, 16)
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("stocks.q", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go ProcessStockData()
	go PublishToQueue(ch)

	router := mux.NewRouter()

	router.HandleFunc("/api/stocks/{stock_code}", func(w http.ResponseWriter, r *http.Request) {
		var body StockRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
		}

		params := mux.Vars(r)
		stockCode, ok := params["stock_code"]
		if !ok {
			http.Error(w, "missing stock_code in request", http.StatusBadRequest)
		}

		iID := xid.New().String()

		reqProcess := RequestProcess{
			InternalID: iID,
			CallerID:   body.ID,
			StockCode:  stockCode,
		}

		fmt.Println("sending process request...", reqProcess)
		Requests <- reqProcess

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}).Methods("POST")

	http.Handle("/", router)
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}
