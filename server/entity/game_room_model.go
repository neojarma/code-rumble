package entity

type GameRoom struct {
	RoomId            string               `json:"roomId"`
	RoomName          string               `json:"roomName"`
	MaxUser           int                  `json:"maxUser"`
	CurrentGameNumber int                  `json:"currentGameNumber"`
	IsGameStart       bool                 `json:"isGameStart"`
	IsRoundStart      bool                 `json:"isRoundStart"`
	IsRoundEnd        bool                 `json:"isRoundEnd"`
	IsGameEnd         bool                 `json:"isGameEnd"`
	Players           []*Player            `json:"players"`
	Question          []*RandQuestion      `json:"-"`
	Leaderboard       []*PlayerLeaderboard `json:"leaderboard"`
}

type NewRoomPayload struct {
	RoomName  string  `json:"roomName"`
	MaxPlayer int     `json:"maxPlayer"`
	Player    *Player `json:"player"`
}

type JoinRoomPayload struct {
	RoomId string  `json:"roomId"`
	Player *Player `json:"player"`
}

type RoomAndQuestion struct {
	Room     *GameRoom `json:"gameRoom"`
	Question *Question `json:"question" gorm:"embedded"`
	StubCode string    `json:"stubCode,omitempty"`
}

type RandQuestion struct {
	Question *Question `json:"question" gorm:"embedded"`
	StubCode string    `json:"stubCode,omitempty"`
}

type PlayerLeaderboard struct {
	PlayerId    string
	DisplayName string
	Score       float64
	Position    int
}
