package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	model "tcp/model"
)

var (
	clients = make(map[net.Conn]model.Message)
	mutex   = sync.Mutex{}
	rooms   = model.Room{}
)

func line(name, role, room, msg string) {
	fmt.Printf("User <%s> with role <%s> send <%s> in room <%s>.\n", name, role, msg, room)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)

	mutex.Lock()
	clients[conn] = model.Message{}
	mutex.Unlock()
	scanner.Scan()
	var m model.Send
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}
		fmt.Println(m.Update)
		if m.Update == true {
			mutex.Lock()
			rooms = m.Message.Room
			mutex.Unlock()
			SendMessage(m.Message, conn, true)
			continue
		}

		line(m.Message.Username, m.Message.Role, m.Message.Room.CurrentRoom, m.Message.Msg)

		mutex.Lock()
		clients[conn] = m.Message
		mutex.Unlock()

		SendMessage(m.Message, conn, false)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
}

func LenUserRoom(m model.Message) (num int) {
	for i := range clients {
		if clients[i].Room.CurrentRoom == m.Room.CurrentRoom {
			num++
		}
	}
	return
}

func SendMessage(m model.Message, conn net.Conn, flag bool) {
	mutex.Lock()
	defer mutex.Unlock()
	var s model.Send
	s.Message.Room.LenUser = LenUserRoom(m)
	s.Message = m
	if flag == true {
		s.Update = true
	}
	data, _ := json.Marshal(s)
	for c, s := range clients {
		if c != conn && s.Room.CurrentRoom == m.Room.CurrentRoom {
			fmt.Fprintln(c, string(data))
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("TCP Server is running on port 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
