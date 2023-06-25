package avalon

import "gaming.candleflame.github.com/gaming/common"

type ShowHistoryHandler struct{}

func (h ShowHistoryHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	return true
}

func (h ShowHistoryHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]

	if alias, ok := gameInfo.playerAddr[playerAddr]; ok {
		for _, msg := range gameInfo.messageHistory[alias] {
			common.SendMessage(playerAddr, msg)
		}
	}
	return "", nil
}
