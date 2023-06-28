package avalon

import (
	"fmt"

	"gaming.candleflame.github.com/gaming/common"
)

type NotifyTaskHandler struct{}

func (h NotifyTaskHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	return gameInfo.gameState == STAGE_NOTIFY_TASK
}

func (h NotifyTaskHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]
	selectedPlayers := gameInfo.selectedPlayer
	roleMap := gameInfo.playerAddr

	boardMsg := fmt.Sprintf("第%v次任务, 你被选为队员执行任务,输入`赞成`完成任务,`反对`破坏任务。",
		len(gameInfo.taskResult)+1)

	for addr, alias := range roleMap {
		for _, selectedAlias := range selectedPlayers {
			if alias == selectedAlias {
				common.SendMessage(addr, boardMsg)
			}
		}
	}

	gameInfo.vote = make(map[string]bool)
	gameInfo.gameState = STAGE_VOTE_FOR_TASK
	return "", nil
}
