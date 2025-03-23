package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	command "tcp/command"
	model "tcp/model"
)

var Info model.Send
var msgTime string
var o bool = false

func myLine(name, role, room string, t bool) {
	var curTime string
	if t {
		curTime = fmt.Sprintf("%s", time.Now().Format("15:04:05"))
		msgTime = curTime
	}
	curTime = msgTime
	fmt.Printf("%s %s@%s:~/%s$ ", curTime, name, role, room)
}

func lineFromServer(name, role, room, msg string) {
	curTime := fmt.Sprintf("%s", time.Now().Format("15:04:05"))
	fmt.Printf("%s %s@%s:~/%s$ %s\n", curTime, name, role, room, msg)
}

func ReadMessage(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var m model.Send
		if err := json.Unmarshal(scanner.Bytes(), &m); err == nil {
			Info.Message.Room = m.Message.Room
			fmt.Println(m.Message.Username, m.Message.Room.Rooms)
			fmt.Println(Info.Message.Username, Info.Message.Room.Rooms)
			if m.Update == true {
				continue
			}
			fmt.Print("\r\033[K")
			lineFromServer(m.Message.Username, m.Message.Role, m.Message.Room.CurrentRoom, m.Message.Msg)
			myLine(Info.Message.Username, Info.Message.Role, Info.Message.Room.CurrentRoom, false)
		}
	}
}

func isCheckMessage(m string) bool {
	if m[0] == 33 {
		return true
	}
	return false
}

func Room() string {
	return "main"
}

func firstConn(scanner *bufio.Scanner, conn net.Conn, m model.Message) (model.Message, error) {
	scanner.Scan()
	m.Username = scanner.Text()
	data, err := json.Marshal(m)
	if err != nil {
		return m, err
	}

	fmt.Fprintln(conn, string(data))
	myLine(m.Username, m.Role, m.Room.CurrentRoom, true)
	return m, nil
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()
	m := model.Message{Role: "user", Room: model.Room{CurrentRoom: Room(), Rooms: []string{"main"}}}
	fmt.Print("Введите ваше имя: ")
	scanner := bufio.NewScanner(os.Stdin)
	m, err = firstConn(scanner, conn, m)
	if err != nil {
		fmt.Println(err)
		return
	}

	go ReadMessage(conn)
	for scanner.Scan() {

		m.Msg = scanner.Text()
		if m.Msg == "" {
			myLine(m.Username, m.Role, m.Room.CurrentRoom, true)
			continue
		}
		if isCheckMessage(m.Msg) {
			Info.Message.Msg = m.Msg
			m, err = command.Command(Info, conn)
			if err != nil {
				fmt.Println(err)
			}
			myLine(m.Username, m.Role, m.Room.CurrentRoom, true)
			Info.Message = m
			Info.Update = true
			data, _ := json.Marshal(Info)
			fmt.Fprintln(conn, string(data))
			Info.Update = false
			continue
		}
		myLine(m.Username, m.Role, m.Room.CurrentRoom, true)
		Info.Message = m
		data, _ := json.Marshal(Info)
		fmt.Fprintln(conn, string(data))
	}
}
