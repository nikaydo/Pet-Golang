package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	c "tcp/client/command"
	"tcp/model"
	"time"
)

var (
	user    model.User
	msgTime string = fmt.Sprintf("%s", time.Now().Format("15:04:05"))
	room    string = user.Room
	rooms          = make(map[string]int)
)

func myLine(t bool) {
	var curTime string
	if t {
		curTime = fmt.Sprintf("%s", time.Now().Format("15:04:05"))
		msgTime = curTime
	}
	curTime = msgTime
	fmt.Printf("%s %s@%s:~/%s$ ", curTime, user.Username, user.Role, room)
}

func lineFromServer(m model.Message) {
	curTime := fmt.Sprintf("%s", time.Now().Format("15:04:05"))
	fmt.Printf("%s %s@%s:~/%s$ %s\n", curTime, m.Username, m.Role, room, m.Msg)
}

func firstConn(conn net.Conn) (u model.User, err error) {
	user.Default()
	usr, _ := json.Marshal(user)
	var event model.Event
	event.Prepare(usr, "user")
	data, _ := json.Marshal(event)
	fmt.Fprintln(conn, string(data))
	return
}

func ReadMessage(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var event model.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			switch event.Type {
			case "user":
				var m model.User
				if err := json.Unmarshal(event.Json, &m); err == nil {
					user = m
					room = m.Room
					fmt.Print("\r\033[K")
					myLine(false)
				}
			case "message":
				var m model.Message
				if err := json.Unmarshal(event.Json, &m); err == nil {
					fmt.Print("\r\033[K")
					lineFromServer(m)
					myLine(false)
				}
			case "room":
				var m model.Command
				if err := json.Unmarshal(event.Json, &m); err == nil {
					fmt.Print("\r\033[K")
					rooms = m.Rooms
					for s, n := range m.Rooms {
						fmt.Println("Room:", s, " Users: ", n)
					}
					myLine(false)
				}
			case "join", "leave", "admin", "admin_err_key", "admin_out":
				Actions(event)
			default:
				continue
			}
		}
	}
}

func Actions(event model.Event) {
	var m model.Command
	if err := json.Unmarshal(event.Json, &m); err == nil {
		fmt.Print("\r\033[K")
		switch event.Type {
		case "join":
			fmt.Println("", m.Username, "joined")
		case "leave":
			fmt.Println("", m.Username, "leave")
		case "admin":
			fmt.Println("Now you are admin!")
		case "admin_out":
			fmt.Println("Now you are user!")
		case "admin_err_key":
			fmt.Println("Wrong key!")
		}
		myLine(false)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()
	var good bool = true
	scanner := bufio.NewScanner(os.Stdin)
	for good {
		fmt.Print("Enter your name: ")
		scanner.Scan()
		n := scanner.Text()
		if n[0] == 33 {
			fmt.Println("Name cannot start with !")
			continue
		}
		user.Username = n
		good = false
	}
	user, err = firstConn(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	go ReadMessage(conn)
	var event model.Event
	var m model.Message = model.Message{ID: user.ID}
	for scanner.Scan() {
		myLine(true)
		m.Msg = scanner.Text()
		if m.Msg == "" {
			continue
		}
		if m.Msg[0] == 33 {
			mod, err := c.Com(m.Msg)

			if err != nil {
				fmt.Println(err)
				continue
			}
			cmd, _ := json.Marshal(mod)
			event.Prepare(cmd, "command")
			data, _ := json.Marshal(event)
			fmt.Fprintln(conn, string(data))
			continue
		}
		msg, _ := json.Marshal(m)
		event.Prepare(msg, "message")
		data, _ := json.Marshal(event)
		fmt.Fprintln(conn, string(data))
	}
}
