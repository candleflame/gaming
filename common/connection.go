package common

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var connections = make(map[string]*websocket.Conn)

func GetConnection(addr string) *websocket.Conn {
	if conn, ok := connections[addr]; ok {
		return conn
	}

	return nil
}

func AddConnection(addr string, conn *websocket.Conn) {
	connections[addr] = conn
}

func DelConnection(addr string) {
	delete(connections, addr)
}

func SendMessage(addr string, msg string) error {
	conn := GetConnection(addr)
	if conn == nil {
		return fmt.Errorf("can not find player")
	}

	// Write response back to client
	if err := conn.WriteMessage(1, []byte(msg)); err != nil {
		return err
	}
	return nil
}
