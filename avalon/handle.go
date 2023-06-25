package avalon

import (
	"fmt"

	"gaming.candleflame.github.com/gaming/common"
)

type GameInfo struct {
	playerNum      int
	playerRoles    map[string]Role
	playerAddr     map[string]string
	gameState      Stage
	taskResult     []bool
	playerOrders   []string
	currentLeader  int
	selectedPlayer []string
	vote           map[string]bool
	messageHistory map[string][]string
}

var gameInfoMap = make(map[string]*GameInfo)

type GameHandler struct {
}

type ActionHandler interface {
	Check(playerAddr string, message *common.ClientMessage) bool
	Handle(playerAddr string, message *common.ClientMessage) (string, error)
}

var actionHandlers = map[Action]ActionHandler{
	ACTION_CREATE_ROOM:   &CreateRoomHandler{},
	ACTION_JOIN:          &JoinHandler{},
	ACTION_INPUT:         &InputHander{},
	ACTION_ASSIGN_ROLE:   &AssignRolesHandler{},
	ACTION_NOTIFY_LEADER: &NotifyLeaderHandler{},
	ACTION_SHOW_HISTORY:  &ShowHistoryHandler{},
}

var AvalonHandler = &GameHandler{}

func (h GameHandler) Handle(removeAddress string, message *common.ClientMessage) (string, error) {
	if handler, ok := actionHandlers[Action(message.Action)]; ok {
		checkResult := handler.Check(removeAddress, message)
		if !checkResult {
			return "", fmt.Errorf("can not operation")
		}
		return handler.Handle(removeAddress, message)
	}
	return "", fmt.Errorf("unsupport operation %v", message.Action)
}
