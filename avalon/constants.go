package avalon

type Action int

const (
	// action from player
	ACTION_CREATE_ROOM Action = 1
	ACTION_JOIN        Action = 2
	ACTION_INPUT       Action = 3

	// action inner
	ACTION_ASSIGN_ROLE   Action = 4
	ACTION_NOTIFY_LEADER Action = 5
	ACTION_SHOW_HISTORY  Action = 6
)

type Stage int

const (
	STAGE_WAITING          Stage = 0
	STAGE_SEND_ROLES       Stage = 1
	STAGE_NOTIFY_LEADER    Stage = 2
	STAGE_SELECT_PEOPLE    Stage = 3
	STAGE_ELECATION        Stage = 4
	STAGE_TASK             Stage = 5
	STAGE_SHOW_TASK_RESULT Stage = 6
	STAGE_END              Stage = 7
)

type Role int

const (
	UNKNOWN         Role = 0
	Merlin          Role = 1
	Assassin        Role = 2
	Percival        Role = 3
	Morgana         Role = 4
	Mordred         Role = 5
	Oberon          Role = 6
	LoyalServant    Role = 7
	MinionOfMordred Role = 8
)

type TaskGroup struct {
	TotalNum int
	EvilNum  int
}

type GameConfig struct {
	Roles      []Role
	TaskGroups []TaskGroup
}

var RoleConfig = map[int]GameConfig{
	5: {
		Roles:      []Role{Merlin, Percival, LoyalServant, Morgana, Assassin},
		TaskGroups: []TaskGroup{{TotalNum: 2, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}, {TotalNum: 2, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}},
	},
	6: {
		Roles:      []Role{Merlin, Percival, LoyalServant, LoyalServant, Morgana, Assassin},
		TaskGroups: []TaskGroup{{TotalNum: 2, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}, {TotalNum: 2, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}},
	},
	7: {
		Roles:      []Role{Merlin, Percival, LoyalServant, LoyalServant, Morgana, Oberon, Assassin},
		TaskGroups: []TaskGroup{{TotalNum: 2, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}, {TotalNum: 3, EvilNum: 1}, {TotalNum: 4, EvilNum: 2}, {TotalNum: 4, EvilNum: 1}},
	},
	8: {
		Roles:      []Role{Merlin, Percival, LoyalServant, LoyalServant, LoyalServant, Morgana, Assassin, MinionOfMordred},
		TaskGroups: []TaskGroup{{TotalNum: 3, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}, {TotalNum: 5, EvilNum: 2}, {TotalNum: 5, EvilNum: 1}},
	},
	9: {
		Roles:      []Role{Merlin, Percival, LoyalServant, LoyalServant, LoyalServant, LoyalServant, Mordred, Morgana, Assassin},
		TaskGroups: []TaskGroup{{TotalNum: 3, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}, {TotalNum: 5, EvilNum: 2}, {TotalNum: 5, EvilNum: 1}},
	},
	10: {
		Roles:      []Role{Merlin, Percival, LoyalServant, LoyalServant, LoyalServant, LoyalServant, Mordred, Morgana, Oberon, Assassin},
		TaskGroups: []TaskGroup{{TotalNum: 3, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}, {TotalNum: 4, EvilNum: 1}, {TotalNum: 5, EvilNum: 2}, {TotalNum: 5, EvilNum: 1}},
	},
}
