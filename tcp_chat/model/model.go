package model

import (
	"encoding/json"
)

type Event struct {
	Type string `json:"type"`
	Json json.RawMessage
}

func (e *Event) Prepare(b json.RawMessage, t string) (data []byte) {
	e.Type = t
	e.Json = b
	data, _ = json.Marshal(e)
	return
}

type User struct {
	ID       int    `json:"id"`
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
	ID       int    `json:"id"`
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

func (c *Command) Prepare(usr User) {

	c.Username = usr.Username
}

func (c *Command) Fill(key, command, action string, rooms map[string]int) {
	c.Key = key
	c.Command = command
	c.Action = action
	c.Rooms = rooms
}

func (m *User) Default() {
	m.Role = "user"
	m.Room = "main"
}
