package avalon

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"gaming.candleflame.github.com/gaming/common"
)

type SelectPeopleHandler struct{}

func (h SelectPeopleHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	alias, ok := gameInfo.playerAddr[playerAddr]
	if !ok {
		return false
	}
	return gameInfo.gameState == STAGE_SELECT_PEOPLE && alias == gameInfo.playerOrders[gameInfo.currentLeader]
}

func (h SelectPeopleHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]

	inputRequest := &InputRequest{}
	json.Unmarshal([]byte(message.Info), inputRequest)
	input := inputRequest.Input
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(input, -1)

	roleList := make([]string, 0)
	for _, match := range matches {
		selectAlias := match[1]
		contains := false
		for _, element := range roleList {
			if element == selectAlias {
				contains = true
				break
			}
		}
		if contains {
			continue
		}
		if _, ok := gameInfo.playerRoles[selectAlias]; !ok {
			continue
		}
		roleList = append(roleList, selectAlias)
	}
	if len(roleList) == RoleConfig[gameInfo.playerNum].TaskGroups[len(gameInfo.taskResult)].TotalNum {

		boardMsg := fmt.Sprintf("%v 选择 ( %v ) 组成小队执行任务", gameInfo.playerAddr[playerAddr], strings.Join(roleList, "; "))
		voteMsg := "如果你同意任务方案，请回复`赞成`，否则回复`反对`"
		for addr, alias := range gameInfo.playerAddr {
			gameInfo.messageHistory[alias] = append(gameInfo.messageHistory[alias], boardMsg)
			common.SendMessage(addr, boardMsg)
			common.SendMessage(addr, voteMsg)
		}

		gameInfo.selectedPlayer = roleList
		gameInfo.vote = make(map[string]bool)
		gameInfo.gameState = STAGE_VOTE_FOR_GROUP
	} else {
		common.SendMessage(playerAddr, "选择小队人数不够，请重新选择。")
	}
	return "", nil
}
