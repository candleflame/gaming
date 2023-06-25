package avalon

import "gaming.candleflame.github.com/gaming/common"

type SelectPeopleHandler struct{}

func (h SelectPeopleHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	return gameInfo.gameState == STAGE_SEND_ROLES
}

func (h SelectPeopleHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	return "", nil
}
