package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	utils "github.com/leonardonicola/chat/pkg"
)

var addr = flag.String("addr", "localhost:3000", "websocket address")

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	// Create client
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("WS(create client): %v", err)
		os.Exit(1)
	}
	defer c.Close()
	comms := make(chan bool)
	go func() {
		defer close(comms)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Fatalf("CLIENT WS(reading): %v", err)
				return
			}
			log.Printf("SERVER: %s", string(msg))
		}
	}()

	for {
		select {
		// if channel closed (on error reading message), break loop
		case <-comms:
			return
		default:
			if err := utils.WriteMessageFromTerminal(c); err != nil {
				return
			}
		}

	}
}
