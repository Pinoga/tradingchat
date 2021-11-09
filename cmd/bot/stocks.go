package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Header = iota
	Values
)

const (
	Symbol = iota
	Date
	Time
	Open
	High
	Low
	Close
	Volume
)

const (
	InternalError = "unexpected error"
	BadRequest    = "stock_code is missing"
	Unavailable   = "service currently unavailable"
)

func ProcessStockData() {
	for data := range Requests {
		fmt.Println("req")

		var queueResponseError error
		var message string

		fmt.Println()
		resp, queueResponseError := requestStockData(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", data.StockCode))
		if queueResponseError == nil {
			message, queueResponseError = parseStockResponse(resp)
			response := StockResponse{
				Message:  message,
				Error:    queueResponseError != nil,
				CallerID: data.CallerID,
			}
			bytes, err := json.Marshal(response)
			if err == nil {
				Responses <- bytes
			}
		}

	}
}

func requestStockData(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(Unavailable)
	}
	return resp, nil
}

func parseStockResponse(resp *http.Response) (string, error) {
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()

	fmt.Println(data)
	if err != nil || len(data) != 2 {
		return "", fmt.Errorf(Unavailable)
	}

	return fmt.Sprintf("%s quote is $%s per share", data[Values][Symbol], data[Values][Close]), nil
}
