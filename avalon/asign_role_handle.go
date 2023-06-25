package avalon

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gaming.candleflame.github.com/gaming/common"
)

type AssignRolesHandler struct{}

func (a AssignRolesHandler) Check(playerAddr string, message *common.ClientMessage) bool {
	gameInfo := gameInfoMap[message.Room]
	return gameInfo.gameState == STAGE_SEND_ROLES
}

func (a AssignRolesHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	gameInfo := gameInfoMap[message.Room]

	assignRoles(gameInfo)
	initMessageHistory(gameInfo)
	sendRoleDescription(gameInfo)
	gameInfo.gameState = STAGE_NOTIFY_LEADER
	AvalonHandler.Handle("", &common.ClientMessage{Game: message.Game, Room: message.Room, Action: int(ACTION_NOTIFY_LEADER)})
	return "", nil
}

func assignRoles(gameInfo *GameInfo) {
	roles := make([]Role, gameInfo.playerNum)
	copy(roles, RoleConfig[gameInfo.playerNum].Roles)

	// 使用当前时间作为种子初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 对数组进行随机排序（洗牌）
	for i := len(roles) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		roles[i], roles[j] = roles[j], roles[i]
	}

	var index = 0
	for _, alias := range gameInfo.playerAddr {
		gameInfo.playerRoles[alias] = roles[index]
		index += 1
	}
}

func initMessageHistory(gameInfo *GameInfo) {
	for alias := range gameInfo.playerRoles {
		gameInfo.messageHistory[alias] = make([]string, 0)
	}
}

func sendRoleDescription(gameInfo *GameInfo) {
	rolesMap := gameInfo.playerAddr

	playerList := make([]string, 0, gameInfo.playerNum)
	for alias := range gameInfo.playerRoles {
		playerList = append(playerList, alias)
	}
	boardMsg := fmt.Sprintf("所有玩家已就绪,玩家顺序为: %v", strings.Join(playerList, ";"))
	gameInfo.playerOrders = playerList

	for addr, alias := range rolesMap {
		gameInfo.messageHistory[alias] = append(gameInfo.messageHistory[alias], boardMsg)
		common.SendMessage(addr, boardMsg)

		role := gameInfo.playerRoles[alias]
		msg := ""
		switch role {
		case Merlin:
			msg += "你的身份是梅林（好人）。"
			msg += "天黑阶段会看见除莫德雷德外的所有坏人。"
			msg += "邪恶玩家有:"
			for otherAlias, otherRole := range gameInfo.playerRoles {
				if otherRole == Morgana || otherRole == Oberon || otherRole == Assassin || otherRole == MinionOfMordred {
					msg += fmt.Sprintf(" %s;", otherAlias)
				}
			}
		case Percival:
			msg += "你的身份是派西维尔（好人）。"
			msg += "天黑阶段会看见梅林和莫甘娜。"
			msg += "你看到的玩家有:"
			for otherAlias, otherRole := range gameInfo.playerRoles {
				if otherRole == Merlin || otherRole == Morgana {
					msg += fmt.Sprintf(" %s;", otherAlias)
				}
			}
		case LoyalServant:
			msg += "你的身份是忠臣（好人）。"
		case Mordred:
			msg += "你的身份是莫德雷德（坏人）。"
			msg += "梅林看不到他。"
			msg += "你看到的邪恶玩家有:"
			for otherAlias, otherRole := range gameInfo.playerRoles {
				if otherRole == Assassin || otherRole == Morgana || otherRole == MinionOfMordred {
					msg += fmt.Sprintf(" %s;", otherAlias)
				}
			}
		case Morgana:
			msg += "你的身份是莫甘娜（坏人）。"
			msg += "假扮梅林，迷惑派西维尔。"
			msg += "你看到的邪恶玩家有:"
			for otherAlias, otherRole := range gameInfo.playerRoles {
				if otherRole == Assassin || otherRole == Mordred || otherRole == MinionOfMordred {
					msg += fmt.Sprintf(" %s;", otherAlias)
				}
			}
		case Oberon:
			msg += "你的身份是奥伯伦（坏人）。"
			msg += "看不到其他坏人，其他坏人也看不到他。"
		case Assassin:
			msg += "你的身份是奥伯伦（坏人）。"
			msg += "在好人阵营3次任务成功后,独自决定,挑选一名可能是梅林的玩家刺杀,如选中,坏人胜利。"
			msg += "你看到的邪恶玩家有:"
			for otherAlias, otherRole := range gameInfo.playerRoles {
				if otherRole == Morgana || otherRole == Mordred || otherRole == MinionOfMordred {
					msg += fmt.Sprintf(" %s;", otherAlias)
				}
			}
		case MinionOfMordred:
			msg += "你的身份是莫德雷德的爪牙（坏人）。"
			msg += "你看到的邪恶玩家有:"
			for otherAlias, otherRole := range gameInfo.playerRoles {
				if otherRole == Morgana || otherRole == Mordred || otherRole == Assassin {
					msg += fmt.Sprintf(" %s;", otherAlias)
				}
			}
		}
		gameInfo.messageHistory[alias] = append(gameInfo.messageHistory[alias], msg)
		common.SendMessage(addr, msg)
	}
}
