package server

import (
	"sync"

	"chat/shared"
)

// MessageStore
// Thread-safe in-memory message queue
type MessageStore struct {
	messages []shared.Message
	nextID   int
	mu       sync.Mutex
}

func NewMessageStore() *MessageStore {
	return &MessageStore{}
}

func (ms *MessageStore) Add(msg shared.Message) shared.Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	msg.ID = ms.nextID
	ms.nextID++
	ms.messages = append(ms.messages, msg)
	return msg
}

func (ms *MessageStore) FetchAfter(lastID int, channel string) []shared.Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var result []shared.Message
	for _, msg := range ms.messages {
		if msg.Channel != channel {
			continue
		}
		if msg.ID > lastID {
			result = append(result, msg)
		}
	}
	return result
}
