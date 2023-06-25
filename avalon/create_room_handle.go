package avalon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"gaming.candleflame.github.com/gaming/common"
)

type CreateRoomRequest struct {
	Num   int    `json:"num"`
	Alias string `json:"alias"`
}

type CreateRoomResponse struct {
	Room string `json:"room"`
}

type CreateRoomHandler struct{}

func (c CreateRoomHandler) Check(playerAddr string, message *common.ClientMessage) bool {

	createRoomInfo := &CreateRoomRequest{}
	err := json.Unmarshal([]byte(message.Info), createRoomInfo)
	if err != nil {
		fmt.Printf("can not parse create room info: %v\n", err)
		return false
	}
	if message.Room == "" {
		fmt.Println("room message is empty")
		return false
	}
	if _, ok := gameInfoMap[message.Room]; ok {
		fmt.Printf("room %v is existed\n", message.Room)
		return false
	}
	if createRoomInfo.Num < 5 || createRoomInfo.Num > 10 {
		fmt.Printf("room player num is error, %v\n", createRoomInfo.Num)
		return false
	}
	if createRoomInfo.Alias == "" {
		fmt.Println("player alias is empty")
		return false
	}
	return true
}

func (c CreateRoomHandler) Handle(playerAddr string, message *common.ClientMessage) (string, error) {
	createRoomInfo := &CreateRoomRequest{}
	_ = json.Unmarshal([]byte(message.Info), createRoomInfo)

	// 使用当前时间作为种子初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	leaderIndex := rand.Intn(createRoomInfo.Num)

	gameInfo := &GameInfo{
		playerNum:   createRoomInfo.Num,
		playerRoles: make(map[string]Role),
		gameState:   STAGE_WAITING,
		playerAddr: map[string]string{
			playerAddr: createRoomInfo.Alias,
		},
		currentLeader:  leaderIndex,
		messageHistory: make(map[string][]string),
		taskResult:     make([]bool, 0, createRoomInfo.Num),
	}

	gameInfoMap[message.Room] = gameInfo

	response, _ := json.Marshal(CreateRoomResponse{
		Room: message.Room,
	})

	common.SendMessage(playerAddr, fmt.Sprintf("您已成功创建房间，当前房间有 %s", createRoomInfo.Alias))
	return string(response), nil
}
