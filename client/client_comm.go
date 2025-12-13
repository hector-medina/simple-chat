package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"chat/shared"
)

const serverURL = "http://localhost:8080/message"

// msg: Message -> send_message()
func SendMessage(msg shared.Message) {
	data, _ := json.Marshal(msg)
	http.Post(serverURL, "application/json", bytes.NewBuffer(data))
}

// check_message()
// Message <-
//
// Polls the server every 10 seconds and calls message_arrived()
func CheckMessages(p *Participant) {
	lastID := -1
	channel := p.Channel()
	for {
		reqURL, err := url.Parse(serverURL)
		if err != nil {
			fmt.Println("invalid server url:", err)
			return
		}
		query := reqURL.Query()
		query.Set("lastId", strconv.Itoa(lastID))
		query.Set("channel", channel)
		reqURL.RawQuery = query.Encode()

		resp, err := http.Get(reqURL.String())
		if err == nil && resp.StatusCode == http.StatusOK {
			var msgs []shared.Message
			json.NewDecoder(resp.Body).Decode(&msgs)
			resp.Body.Close()

			for _, msg := range msgs {
				p.MessageArrived(msg)
				if msg.ID > lastID {
					lastID = msg.ID
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
