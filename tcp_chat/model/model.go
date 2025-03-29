package model

import (
	"encoding/json"
)

type Event struct {
	Type string          `json:"type"`
	Json json.RawMessage `json:"json"`
}

func (e *Event) Prepare(b json.RawMessage, t string) (data string) {
	e.Type = t
	e.Json = b
	d, _ := json.Marshal(e)
	data = string(d)
	return
}

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Room     string `json:"rooms"`
}

func (m *User) Prepare(u, r, room string) {
	m.Username = u
	m.Role = r
	m.Role = room
}

type Message struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Msg      string `json:"message"`
}

func (m *Message) Prepare(u, r, msg string) {
	m.Username = u
	m.Role = r
	m.Msg = msg
}

type Command struct {
	Key      string         `json:"key"`
	Command  string         `json:"command"`
	Action   string         `json:"user"`
	Username string         `json:"username"`
	Rooms    map[string]int `json:"rooms"`
}

func (m *User) Default() {
	m.Role = "user"
	m.Room = "main"
}
