package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendCommandToBot(c string, bg int) error {
	postBody, _ := json.Marshal(struct {
		ID string `json:"ID"`
	}{
		ID: fmt.Sprint(bg),
	})

	if c == "" {
		return fmt.Errorf("please input a valid command")
	}

	resp, err := http.Post("http://127.0.0.1:9000/api/stocks/"+c, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return fmt.Errorf("unexpected error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(string(body))
	}

	return nil
}
