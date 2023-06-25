package main

import (
	"encoding/json"
	"fmt"

	"gaming.candleflame.github.com/gaming/avalon"
	"gaming.candleflame.github.com/gaming/common"
)

type GameType int

const avalonType GameType = 1

var gameHandlerMap = map[GameType]common.GameHandler{
	avalonType: avalon.AvalonHandler,
}

func Handle(addr string, p []byte) common.ServerMessage {
	clientMessage := &common.ClientMessage{}
	err := json.Unmarshal(p, clientMessage)
	if err != nil {
		fmt.Printf("JSON 解码失败:%v, data %s, client %s", err, p, addr)
		return common.ServerMessage{Status: -1, Msg: "请求内容解码失败"}
	}

	if handler, ok := gameHandlerMap[GameType(clientMessage.Game)]; ok {
		message, err := handler.Handle(addr, clientMessage)
		if err != nil {
			return common.ServerMessage{Status: -1, Msg: fmt.Sprintf("游戏处理异常, %v", err)}
		}
		return common.ServerMessage{Status: 0, Data: message}
	} else {
		return common.ServerMessage{Status: -1, Msg: fmt.Sprintf("游戏不存在, %v", clientMessage.Game)}
	}
}
