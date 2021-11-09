package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

func SendJSONResponse(w http.ResponseWriter, message string, status int, data interface{}) {
	resp := Response{}

	resp.Message = message
	resp.Status = fmt.Sprint(status)
	resp.Data = data

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("could not marshal response. Err: %v", err)
	}

	w.WriteHeader(status)
	w.Write(jsonResp)
}
