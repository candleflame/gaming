package avalon

import (
	"fmt"

	"gaming.candleflame.github.com/gaming/common"
)

type NotifyLeaderHandler struct{}

func (n NotifyLeaderHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	return gameInfo.gameState == STAGE_NOTIFY_LEADER
}

func (n NotifyLeaderHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]
	rolesMap := gameInfo.playerAddr

	gameInfo.currentLeader = (gameInfo.currentLeader + 1) % gameInfo.playerNum
	boardMsg := fmt.Sprintf("第%v次任务,由 %v 来选择任务成员，需要选择 %v 人组成小队，若其中有 %v 人破坏任务，则任务失败",
		len(gameInfo.taskResult)+1,
		gameInfo.playerOrders[gameInfo.currentLeader],
		RoleConfig[gameInfo.playerNum].TaskGroups[len(gameInfo.taskResult)].TotalNum,
		RoleConfig[gameInfo.playerNum].TaskGroups[len(gameInfo.taskResult)].EvilNum)
	for addr, alias := range rolesMap {
		gameInfo.messageHistory[alias] = append(gameInfo.messageHistory[alias], boardMsg)
		common.SendMessage(addr, boardMsg)
		if alias == gameInfo.playerOrders[gameInfo.currentLeader] {
			common.SendMessage(addr, "点击角色按钮，即为选择成为队友")
		}
	}

	gameInfo.gameState = STAGE_SELECT_PEOPLE
	return "", nil
}
