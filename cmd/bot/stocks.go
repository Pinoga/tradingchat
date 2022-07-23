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
	BadRequest    = "stock_code is invalid"
	Unavailable   = "service currently unavailable"
)

func ProcessStockData() {
	for data := range Requests {
		fmt.Println("req")

		var queueResponseError error

		resp, queueResponseError := requestStockData(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", data.StockCode))
		if queueResponseError == nil {
			message, hasErr := parseStockResponse(resp)
			response := StockResponse{
				Message:  message,
				Error:    hasErr,
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

func parseStockResponse(resp *http.Response) (string, bool) {
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()

	fmt.Println(data)
	if err != nil || len(data) != 2 {
		return Unavailable, true
	}
	if data[Values][Date] == "N/D" {
		return BadRequest, true
	}

	return fmt.Sprintf("%s quote is $%s per share", data[Values][Symbol], data[Values][Close]), false
}
