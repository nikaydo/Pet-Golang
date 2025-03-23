package model

type Message struct {
	Username string `json:"username"`
	Msg      string `json:"message"`
	Role     string `json:"role"`
	Room     Room   `json:"room"`
}

type Send struct {
	Message Message `json:"message"`
	Update  bool    `json:"update"`
}
type Room struct {
	CurrentRoom string   `json:"currentRoom"`
	LenUser     int      `json:"lenUser"`
	Rooms       []string `json:"rooms"`
}
