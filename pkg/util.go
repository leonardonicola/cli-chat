package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func WriteMessageFromTerminal(conn *websocket.Conn) error {

	// creater a reader to read from the terminal
	reader := bufio.NewReader(os.Stdin)

	// read strings with new line delimiter
	msg, _ := reader.ReadString('\n')

	// write to connection
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("WS(write): %v", err)
		return fmt.Errorf("WS(write): %v", err)
	}
	return nil
}
