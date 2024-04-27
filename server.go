package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	utils "github.com/leonardonicola/chat/pkg"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const PORT = ":3000"

func main() {
	var addr = flag.String("addr", fmt.Sprintf("localhost%s", PORT), "http address")
	url := url.URL{Scheme: "http", Host: *addr, Path: "/"}
	http.HandleFunc(url.Path, handler)

	forever := make(chan bool)
	go func() {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatalf("HTTP(listen): %v", err)
			forever <- false
		}
	}()
	log.Printf("LISTENING AT: %s", url.String())
	<-forever

}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("WS(upgrade): %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WS(reading message): %v", err)
				return
			}
			log.Printf("CLIENT: %s", string(msg))
		}
	}()

	for {
		select {
		// if channel closed (on error reading message), break loop
		case <-done:
			return
		default:
			if err := utils.WriteMessageFromTerminal(conn); err != nil {
				return
			}
		}

	}
}
