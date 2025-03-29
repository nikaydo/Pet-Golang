package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"tcp/model"
	roomHandler "tcp/server/handlers"
)

var (
	clients = make(map[net.Conn]model.User)
	mutex   = sync.Mutex{}
	rooms   = map[string]int{"main": 0}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())
	var m model.User
	mutex.Lock()
	clients[conn] = m
	rooms["main"]++
	mutex.Unlock()
	scanner := bufio.NewScanner(conn)
	var event model.Event
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}
		switch event.Type {
		case "message":
			var m model.Message
			if err := json.Unmarshal(event.Json, &m); err == nil {
				Line(conn, m)
				m.Prepare(clients[conn].Username, clients[conn].Role, m.Msg)
				msg, err := json.Marshal(m)
				if err != nil {
					fmt.Println(err)
					continue
				}
				SendAction(event.Prepare(msg, "message"), clients[conn], conn)
			}
		case "user":
			var u model.User = model.User{}
			if err := json.Unmarshal(event.Json, &u); err == nil {
				mutex.Lock()
				clients[conn] = u
				mutex.Unlock()
				SendUser(conn)
				SendCommand(event, conn, clients[conn].Username+" joined")
			}
		case "command":
			var c model.Command
			if err := json.Unmarshal(event.Json, &c); err == nil {
				switch c.Command {
				case "admin":
					admin(c, conn, event)
				case "room":
					Room(c, conn, rooms)
				}
			}
		}
	}

	SendCommand(event, conn, clients[conn].Username+" disconnected")
	mutex.Lock()
	rooms[clients[conn].Room]--
	delete(clients, conn)
	mutex.Unlock()
}

func SendCommand(event model.Event, conn net.Conn, t string) {
	var usrEvent model.Command
	usrEvent.Username = clients[conn].Username
	usr, err := json.Marshal(usrEvent)
	if err != nil {
		fmt.Println(err)
		return
	}
	SendAction(event.Prepare(usr, t), clients[conn], conn)
}

func Room(c model.Command, conn net.Conn, rooms map[string]int) {
	var event model.Event
	switch c.Action {
	case "list":
		c.Rooms = rooms
		roomHandler.List(event, c, conn)
	case "join":
		oldRoom := clients[conn].Room
		for i, _ := range rooms {
			if i == c.Key {
				rooms[oldRoom]--
				mutex.Lock()
				n := clients[conn]
				n.Room = c.Key
				clients[conn] = n
				mutex.Unlock()
				rooms[c.Key]++
				return
			}
		}
		SendUser(conn)
	case "add":
		roomHandler.Add(event, c, conn, rooms)
		SendUser(conn)
	case "del":
		delete(rooms, c.Key)
		c.Rooms = rooms
		roomHandler.Del(event, c, conn)
		SendUser(conn)
	case "info":
		c.Rooms = rooms
		c.Action = "Room : " + c.Key + " Users : " + fmt.Sprint(rooms[c.Key])
		roomHandler.Info(event, c, conn)
	}
}

func adminResponce(conn net.Conn, event model.Event, role string, t string) {
	mutex.Lock()
	n := clients[conn]
	n.Role = role
	clients[conn] = n
	mutex.Unlock()
	SendCommand(event, conn, t)
	SendUser(conn)
}

func admin(c model.Command, conn net.Conn, event model.Event) {
	switch c.Key {
	case "admin":
		adminResponce(conn, event, "admin", "You are admin")
	case "out":
		adminResponce(conn, event, "user", "You are user")
	case "admin_err_key":
		SendCommand(event, conn, "Error key")
	}
}

func SendUser(conn net.Conn) {
	var event model.Event
	usr, err := json.Marshal(clients[conn])
	if err != nil {
		fmt.Println(err)
		return
	}
	data := event.Prepare(usr, "user")
	fmt.Fprintln(conn, (data))
	return
}

func SendAction(data string, usr model.User, conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for c, s := range clients {
		if c != conn && s.Room == usr.Room {
			fmt.Fprintln(c, data)
		}
	}
}

func Line(conn net.Conn, m model.Message) {
	usr := clients[conn]
	fmt.Printf("User <%s> with role <%s> send <%s> in room <%s>.\n", usr.Username, usr.Role, m.Msg, usr.Room)
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
