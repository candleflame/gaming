package avalon

import (
	"encoding/json"
	"fmt"
	"strings"

	"gaming.candleflame.github.com/gaming/common"
)

type VoteHandler struct{}

func (h VoteHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	if gameInfo.gameState == STAGE_VOTE_FOR_GROUP {
		return true
	}

	if gameInfo.gameState == STAGE_VOTE_FOR_TASK {
		alias, ok := gameInfo.playerAddr[playerAddr]
		if !ok {
			return false
		}
		playerList := gameInfo.selectedPlayer
		for _, player := range playerList {
			if player == alias {
				return true
			}
		}
	}
	return false
}

func (h VoteHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]
	alias := gameInfo.playerAddr[playerAddr]

	inputRequest := &InputRequest{}
	json.Unmarshal([]byte(message.Info), inputRequest)
	isAgree := inputRequest.Input == "赞成"

	common.SendMessage(playerAddr, "您已完成投票，等待其他玩家投票。")

	if gameInfo.gameState == STAGE_VOTE_FOR_GROUP {
		gameInfo.vote[alias] = isAgree

		if len(gameInfo.vote) == gameInfo.playerNum {
			agreeNum := 0
			for _, isAgree := range gameInfo.vote {
				if isAgree {
					agreeNum += 1
				}
			}

			if agreeNum > gameInfo.playerNum/2 {
				gameInfo.gameState = STAGE_NOTIFY_TASK
				AvalonHandler.Handle("", &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_NOTIFY_TASK)})
			} else {
				gameInfo.gameState = STAGE_NOTIFY_LEADER
				AvalonHandler.Handle("", &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_NOTIFY_LEADER)})
			}
		}
		return fmt.Sprintf("%v 投票完成", alias), nil
	}

	if gameInfo.gameState == STAGE_VOTE_FOR_TASK {
		gameInfo.vote[alias] = isAgree
		if len(gameInfo.vote) == len(gameInfo.selectedPlayer) {
			disagreeNum := 0
			for _, isAgree := range gameInfo.vote {
				if !isAgree {
					disagreeNum += 1
				}
			}

			taskResult := disagreeNum < RoleConfig[gameInfo.playerNum].TaskGroups[len(gameInfo.taskResult)].EvilNum
			gameInfo.taskResult = append(gameInfo.taskResult, taskResult)
			taskResultStr := "失败"
			if taskResult {
				taskResultStr = "成功"
			}
			boardMsg := fmt.Sprintf("第%v次任务, 执行任务的人员有：%v, 其中反对票有 %v, 任务%v。",
				len(gameInfo.taskResult), strings.Join(gameInfo.selectedPlayer, ";"), disagreeNum, taskResultStr)
			for addr, alias := range gameInfo.playerAddr {
				gameInfo.messageHistory[alias] = append(gameInfo.messageHistory[alias], boardMsg)
				common.SendMessage(addr, boardMsg)
			}

			gameInfo.gameState = STAGE_JUDGE_RESULT
			AvalonHandler.Handle("", &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_JUDGE_RESULT)})
		}

	}
	return "", nil
}
