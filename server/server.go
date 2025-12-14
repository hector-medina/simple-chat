package server

import (
	"encoding/json"
	"log"

	"chat/shared"

	zmq "github.com/pebbe/zmq4"
)

const (
	pullBindAddress = "tcp://*:5555"
	pubBindAddress  = "tcp://*:5563"
)

func StartServer() {
	store := NewMessageStore()
	chat := NewChat(store)

	pull, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		log.Fatalf("unable to create pull socket: %v", err)
	}
	defer pull.Close()
	if err := pull.Bind(pullBindAddress); err != nil {
		log.Fatalf("unable to bind pull socket: %v", err)
	}

	pub, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatalf("unable to create pub socket: %v", err)
	}
	defer pub.Close()
	if err := pub.Bind(pubBindAddress); err != nil {
		log.Fatalf("unable to bind pub socket: %v", err)
	}

	log.Printf("ZeroMQ server ready (push-bind=%s pub-bind=%s)", pullBindAddress, pubBindAddress)

	for {
		raw, err := pull.RecvBytes(0)
		if err != nil {
			log.Printf("failed to receive message: %v", err)
			continue
		}

		var msg shared.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("invalid payload: %v", err)
			continue
		}

		stored := chat.DistributeMessage(msg)

		payload, err := json.Marshal(stored)
		if err != nil {
			log.Printf("unable to marshal message for broadcast: %v", err)
			continue
		}

		if _, err := pub.Send(stored.Channel, zmq.SNDMORE); err != nil {
			log.Printf("unable to send channel envelope: %v", err)
			continue
		}

		if _, err := pub.SendBytes(payload, 0); err != nil {
			log.Printf("unable to broadcast message: %v", err)
		}
	}
}
