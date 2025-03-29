package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"tcp/model"
)

func List(event model.Event, c model.Command, conn net.Conn) {
	r, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(conn, (event.Prepare(r, "room")))
}
func Add(event model.Event, c model.Command, conn net.Conn, rooms map[string]int) {
	for i, _ := range rooms {
		if i == c.Key {
			r, err := json.Marshal(c)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Fprintln(conn, (event.Prepare(r, "room already exists")))
			return
		}
	}
	rooms[c.Key] = 0
	c.Rooms = rooms
	r, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(conn, (event.Prepare(r, "room add")))
}

func Del(event model.Event, c model.Command, conn net.Conn) {
	r, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(conn, (event.Prepare(r, "room delete")))

}

func Info(event model.Event, c model.Command, conn net.Conn) {
	r, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(conn, (event.Prepare(r, "room")))
}

