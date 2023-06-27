package avalon

import (
	"encoding/json"
	"fmt"
	"regexp"

	"gaming.candleflame.github.com/gaming/common"
)

type InputType int

const (
	HISTORY_QUERY InputType = 1
	SELECT_GROUP  InputType = 2
	DEFAULT_QUERY InputType = -1
)

var inputParseArr = []InputType{HISTORY_QUERY, SELECT_GROUP, DEFAULT_QUERY}

type InputTypeHandle struct {
	Reg    *regexp.Regexp
	Handle func(string, *GameInfo, *common.ClientMessage) (string, error)
}

var inputTypeHandleMap = map[InputType]InputTypeHandle{
	HISTORY_QUERY: {Reg: regexp.MustCompile("^历史$"), Handle: HistoryQueryHandle},
	SELECT_GROUP:  {Reg: regexp.MustCompile("^选择.*"), Handle: SelectGroupHandle},
	DEFAULT_QUERY: {Reg: regexp.MustCompile(".*"), Handle: DefaultHandle},
}

type InputRequest struct {
	Input string `json:"input"`
}

type InputHander struct{}

func (i InputHander) Check(playerAddr string, message *common.ClientMessage) bool {
	inputRequest := &InputRequest{}
	err := json.Unmarshal([]byte(message.Info), inputRequest)
	if err != nil {
		fmt.Printf("can not parse input info: %v\n", err)
		return false
	}
	if inputRequest.Input == "" {
		fmt.Printf("input info is empty\n")
		return false
	}
	return true
}

func (i InputHander) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	inputRequest := &InputRequest{}
	json.Unmarshal([]byte(message.Info), inputRequest)
	input := inputRequest.Input
	gameInfo := gameInfoMap[message.Room]

	for _, inputType := range inputParseArr {
		inputTypeHandle := inputTypeHandleMap[inputType]
		if inputTypeHandle.Reg.MatchString(input) {
			inputTypeHandle.Handle(playerAddr, gameInfo, message)
			break
		}
	}
	return "", nil
}

func HistoryQueryHandle(playerAddr string, gameInfo *GameInfo, message *common.ClientMessage) (string, error) {
	return AvalonHandler.Handle(playerAddr, &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_SHOW_HISTORY)})
}

func SelectGroupHandle(playerAddr string, gameInfo *GameInfo, message *common.ClientMessage) (string, error) {
	return AvalonHandler.Handle(playerAddr, &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_SELECT_GROUP), Info: message.Info})
}

func DefaultHandle(playerAddr string, gameInfo *GameInfo, message *common.ClientMessage) (string, error) {
	msg := `发送消息规则:\n
	1. 发送“历史” 获取到历史消息。\n
	`
	common.SendMessage(playerAddr, msg)
	return "", nil
}
