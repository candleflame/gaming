package avalon

import (
	"gaming.candleflame.github.com/gaming/common"
)

type JudgeResultHandler struct{}

func (h JudgeResultHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	return gameInfo.gameState == STAGE_JUDGE_RESULT
}

func (h JudgeResultHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]
	taskResult := gameInfo.taskResult
	roleMap := gameInfo.playerAddr

	successNum, failNum := 0, 0
	for _, result := range taskResult {
		if result {
			successNum += 1
		} else {
			failNum += 1
		}
	}

	if failNum >= 3 {
		boardMsg := "坏人成功破坏任务三次,坏人胜利。"
		gameInfo.gameState = STAGE_END
		for addr, _ := range roleMap {
			common.SendMessage(addr, boardMsg)
		}
		return "", nil
	}

	if successNum >= 3 {
		boardMsg := "好人成功完成任务三次,等待刺客刺杀梅林。"
		gameInfo.gameState = STAGE_KILL_MERLIN
		for addr, alias := range roleMap {
			common.SendMessage(addr, boardMsg)
			role := gameInfo.playerRoles[alias]
			if role == Assassin {
				common.SendMessage(addr, "选择你认为是梅林的玩家进行刺杀")
			}
		}
		return "", nil
	}

	gameInfo.gameState = STAGE_NOTIFY_LEADER
	AvalonHandler.Handle("", &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_NOTIFY_LEADER)})

	return "", nil
}
