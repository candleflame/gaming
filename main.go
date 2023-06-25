package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gaming.candleflame.github.com/gaming/common"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var addr = conn.RemoteAddr().String()
	common.AddConnection(addr, conn)

	defer conn.Close()
	defer common.DelConnection(addr)

	for {
		// Read message from client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.CloseMessage {
			common.DelConnection(addr)
		}
		if string(p) == "ping" {
			// Write response back to client
			if err := conn.WriteMessage(1, []byte("pong")); err != nil {
				return
			}
			continue
		}

		msg := Handle(addr, p)
		response, _ := json.Marshal(msg)
		log.Println(string(response))
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
