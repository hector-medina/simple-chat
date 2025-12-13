package client

import (
	"fmt"
	"sync"

	"chat/shared"
)

// Participant
// Logical design:
//
// nick: Text -> constructor() ->
// userInput: Text -> text_read() ->
// msg: Message -> message_arrived() ->
type Participant struct {
	nick    string
	channel string
	mu      sync.Mutex
}

// nick: Text -> constructor() ->
func NewParticipant(nick, channel string) *Participant {
	return &Participant{
		nick:    nick,
		channel: channel,
	}
}

// Channel returns participant's current channel.
func (p *Participant) Channel() string {
	return p.channel
}

// userInput: Text -> text_read() ->
// Sends the user-written text to the server
func (p *Participant) TextRead(text string) {
	msg := shared.Message{
		Channel: p.channel,
		Author:  p.nick,
		Message: text,
	}
	SendMessage(msg)
}

// msg: Message -> message_arrived() ->
// Displays an incoming message on screen
func (p *Participant) MessageArrived(msg shared.Message) {
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("[%s][%s]: %s\n", msg.Channel, msg.Author, msg.Message)
}
