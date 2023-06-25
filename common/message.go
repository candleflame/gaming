package common

type ClientMessage struct {
	Game   int    `json:"game"`
	Room   string `json:"room"`
	Action int    `json:"action"`
	Info   string `json:"info"`
}

type ServerMessage struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

type GameHandler interface {
	Handle(removeAddress string, message *ClientMessage) (string, error)
}
