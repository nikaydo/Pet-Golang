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
)

func myLine(t bool) {
	var curTime string
	if t {
		curTime = fmt.Sprintf("%s", time.Now().Format("15:04:05"))
		msgTime = curTime
	}
	curTime = msgTime
	fmt.Printf("%s %s@%s:~/%s$ ", curTime, user.Username, user.Role, user.Room)
}

func lineFromServer(m model.Message) {
	curTime := fmt.Sprintf("%s", time.Now().Format("15:04:05"))
	fmt.Printf("%s %s@%s:~/%s$ %s\n", curTime, m.Username, m.Role, user.Room, m.Msg)
}

func firstConn(conn net.Conn) (u model.User) {
	user.Default()
	usr, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	var event model.Event
	data := event.Prepare(usr, "user")
	_, err = fmt.Fprintln(conn, data)
	if err != nil {
		fmt.Println("Error:", err)
	}
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
					for s, n := range m.Rooms {
						fmt.Println("Room:", s, " Users: ", n)
					}
					myLine(false)
				}
			default:
				var m model.Command
				if err := json.Unmarshal(event.Json, &m); err == nil {
					fmt.Print("\r\033[K")
					fmt.Println(event.Type)
					myLine(false)
				}
			}
		}
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
	user = firstConn(conn)
	go ReadMessage(conn)
	var event model.Event
	var m model.Message
	for scanner.Scan() {
		myLine(true)
		m.Msg = scanner.Text()
		if m.Msg == "" {
			continue
		}
		if m.Msg[0] == 33 {
			fmt.Print("\r\033[K")
			mod, err := c.Com(m.Msg)
			if err != nil {
				fmt.Println(err)
				myLine(true)
				continue
			}
			cmd, err := json.Marshal(mod)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Fprintln(conn, event.Prepare(cmd, "command"))
			continue
		}
		msg, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)

		}
		fmt.Fprintln(conn, event.Prepare(msg, "message"))
	}
}
