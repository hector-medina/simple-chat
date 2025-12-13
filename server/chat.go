package server

import (
	"log"

	"chat/shared"
)

// Chat
// Logical design:
//
// msg: Message distribute_message() ->
type Chat struct {
	store *MessageStore
}

func NewChat(store *MessageStore) *Chat {
	return &Chat{store: store}
}

// msg: Message distribute_message() ->
// Stores message to be delivered
func (c *Chat) DistributeMessage(msg shared.Message) {
	log.Printf("channel=%s author=%s msg=%q", msg.Channel, msg.Author, msg.Message)
	c.store.Add(msg)
}
