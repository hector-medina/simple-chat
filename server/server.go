package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"chat/shared"
)

func StartServer() {
	store := NewMessageStore()
	chat := NewChat(store)

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		// POST /message
		case http.MethodPost:
			var msg shared.Message
			json.NewDecoder(r.Body).Decode(&msg)
			chat.DistributeMessage(msg)
			w.WriteHeader(http.StatusCreated)

			// GET /message
		case http.MethodGet:
			lastID := -1
			if lastIDParam := r.URL.Query().Get("lastId"); lastIDParam != "" {
				if parsed, err := strconv.Atoi(lastIDParam); err == nil {
					lastID = parsed
				}
			}
			channel := r.URL.Query().Get("channel")
			if channel == "" {
				channel = "general"
			}
			msgs := store.FetchAfter(lastID, channel)
			json.NewEncoder(w).Encode(msgs)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
