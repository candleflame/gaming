package avalon

import (
	"encoding/json"
	"fmt"

	"gaming.candleflame.github.com/gaming/common"
)

type JoinRequest struct {
	Alias string `json:"alias"`
}

type JoinHandler struct{}

func (c JoinHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	joinRequest := &JoinRequest{}
	err := json.Unmarshal([]byte(message.Info), joinRequest)
	if err != nil {
		fmt.Printf("can not parse join room info: %v\n", err)
		return false
	}
	if joinRequest.Alias == "" {
		fmt.Println("player alias is empty")
		return false
	}
	var gameInfo *GameInfo
	var ok bool
	if gameInfo, ok = gameInfoMap[message.Room]; !ok {
		fmt.Printf("room %v is not existed", message.Room)
		return false
	}
	if alias, ok := gameInfo.playerAddr[playerAddr]; ok && alias == joinRequest.Alias {
		return false
	}
	if len(gameInfo.playerAddr) == gameInfo.playerNum {
		isSamePlayer := false
		for _, alias := range gameInfo.playerAddr {
			if alias == joinRequest.Alias {
				isSamePlayer = true
			}
		}
		if !isSamePlayer {
			fmt.Printf("room is full")
			return false
		}
	}
	return true
}

func (c JoinHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	joinRequest := &JoinRequest{}
	_ = json.Unmarshal([]byte(message.Info), joinRequest)

	gameInfo := gameInfoMap[message.Room]
	needToDelete := ""
	for addr, alias := range gameInfo.playerAddr {
		if alias == joinRequest.Alias {
			needToDelete = addr
		}
	}
	if needToDelete != "" {
		delete(gameInfo.playerAddr, needToDelete)
	}

	gameInfo.playerAddr[playerAddr] = joinRequest.Alias
	var response = fmt.Sprintf("欢迎 %v 加入房间，当前房间有: ", joinRequest.Alias)

	for _, alias := range gameInfo.playerAddr {
		response += alias + ", "
	}

	for addr := range gameInfo.playerAddr {
		common.SendMessage(addr, response)
	}

	if gameInfo.gameState == STAGE_WAITING && gameInfo.playerNum == len(gameInfo.playerAddr) {
		gameInfo.gameState = STAGE_SEND_ROLES
		AvalonHandler.Handle("", &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_ASSIGN_ROLE)})

	}

	return fmt.Sprintf("%s join room %s successful", joinRequest.Alias, message.Room), nil
}
