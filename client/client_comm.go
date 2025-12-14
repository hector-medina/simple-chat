package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"chat/shared"

	zmq "github.com/pebbe/zmq4"
)

const (
	defaultServerHost = "localhost"
)

var (
	serverHost        = getServerHost()
	pushServerAddress = fmt.Sprintf("tcp://%s:5555", serverHost)
	pubServerAddress  = fmt.Sprintf("tcp://%s:5563", serverHost)
)

func getServerHost() string {
	if host := os.Getenv("CHAT_SERVER_HOST"); host != "" {
		return host
	}
	return defaultServerHost
}

var (
	pushSocket     *zmq.Socket
	pushSocketOnce sync.Once
	pushSocketErr  error
)

func getPushSocket() (*zmq.Socket, error) {
	pushSocketOnce.Do(func() {
		var err error
		pushSocket, err = zmq.NewSocket(zmq.PUSH)
		if err != nil {
			pushSocketErr = err
			return
		}
		if err = pushSocket.Connect(pushServerAddress); err != nil {
			pushSocketErr = err
		}
	})
	return pushSocket, pushSocketErr
}

// msg: Message -> send_message()
func SendMessage(msg shared.Message) {
	socket, err := getPushSocket()
	if err != nil {
		log.Printf("unable to get push socket: %v", err)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("unable to marshal message: %v", err)
		return
	}

	if _, err := socket.SendBytes(data, 0); err != nil {
		log.Printf("unable to send message: %v", err)
	}
}

// check_message()
// Message <-
//
// Subscribes to the server and calls message_arrived()
func CheckMessages(p *Participant) {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Printf("unable to create subscriber: %v", err)
		return
	}
	defer subscriber.Close()

	if err := subscriber.Connect(pubServerAddress); err != nil {
		log.Printf("unable to connect subscriber: %v", err)
		return
	}

	if err := subscriber.SetSubscribe(p.Channel()); err != nil {
		log.Printf("unable to subscribe to channel %s: %v", p.Channel(), err)
		return
	}

	for {
		// envelope (channel)
		if _, err := subscriber.Recv(0); err != nil {
			log.Printf("unable to receive envelope: %v", err)
			continue
		}

		payload, err := subscriber.RecvBytes(0)
		if err != nil {
			log.Printf("unable to receive payload: %v", err)
			continue
		}

		var msg shared.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			log.Printf("unable to decode payload: %v", err)
			continue
		}

		p.MessageArrived(msg)
	}
}
