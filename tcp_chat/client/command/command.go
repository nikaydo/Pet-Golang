package command

import (
	"fmt"
	"strings"
	"tcp/model"
)

var (
	ErrNoParam = fmt.Errorf("no param")
	ErrUnknown = fmt.Errorf("unknown command")
)

func Com(msg string) (c model.Command, err error) {
	msg = strings.Replace(msg, "!", "", 1)
	list := strings.Split(msg, " ")
	switch list[0] {
	case "help":
		cHelp()
	case "admin":
		f := Admin(list)
		c.Key = f
		c.Command = list[0]
	case "room":
		c = Room(list)
	default:
		return c, ErrUnknown
	}
	return
}

func Room(list []string) (c model.Command) {
	c.Command = "room"
	switch list[1] {
	case "list":
		c.Command = "room"
		c.Action = "list"
	case "add":
		c.Command = "room"
		c.Action = "add"
		c.Key = list[2]
	case "del":
		c.Command = "room"
		c.Action = "del"
		c.Key = list[2]
	case "info":
		c.Command = "room"
		c.Action = "info"
		c.Key = list[2]
	case "join":
		c.Command = "room"
		c.Action = "join"
		c.Key = list[2]
	}
	return
}

func cHelp() {
	fmt.Println("help - list of commands")
	fmt.Println("info - about user")
	fmt.Println("room - work with rooms\n    [name] - go to room\n    add  - make new room\n    del  - delete room\n    list - list of rooms\n    info - info about room")
	fmt.Println("admin - admin mode\n    get [key]  - be admin\n    out  - quit from admin mode")
}

func Admin(list []string) (key string) {
	switch list[1] {
	case "get":
		if len(list) < 3 {
			fmt.Println("Error:", ErrNoParam)
			return
		}
		if list[2] == "admin" {
			key = list[2]
		}
	case "out":
		key = "out"
	}
	return
}
