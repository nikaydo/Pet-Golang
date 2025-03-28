package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"tcp/model"
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
				event.Type = "message"
				m.Username = clients[conn].Username
				m.Role = clients[conn].Role
				msg, _ := json.Marshal(m)
				event.Json = msg
				SendAction(event, clients[conn], conn)
			}
		case "user":
			var u model.User = model.User{ID: MakeId()}
			if err := json.Unmarshal(event.Json, &u); err == nil {
				mutex.Lock()
				clients[conn] = u
				mutex.Unlock()
				SendUser(conn, u.ID)
				var usrEvent model.Command
				usrEvent.Username = clients[conn].Username
				usr, _ := json.Marshal(usrEvent)
				event.Prepare(usr, "join")
				SendAction(event, clients[conn], conn)
			}
		case "command":
			var c model.Command
			if err := json.Unmarshal(event.Json, &c); err == nil {
				switch c.Command {
				case "admin":
					admin(c, conn, event)
				case "room":
					Room(c, conn)
				}
			}
		}
	}
	var usrEvent model.Command
	usrEvent.Username = clients[conn].Username
	usr, _ := json.Marshal(usrEvent)
	event.Prepare(usr, "leave")
	SendAction(event, clients[conn], conn)
	mutex.Lock()
	rooms[clients[conn].Room]++
	delete(clients, conn)
	mutex.Unlock()
}

func SendCommand(event model.Event, conn net.Conn, t string) {
	var usrEvent model.Command
	usrEvent.Username = clients[conn].Username
	usr, _ := json.Marshal(usrEvent)
	event.Type = t
	event.Json = usr
	data, _ := json.Marshal(event)
	fmt.Fprintln(conn, string(data))
}

func Room(c model.Command, conn net.Conn) {
	switch c.Action {
	case "list":
		var event model.Event
		c.Rooms = rooms
		r, _ := json.Marshal(c)
		event.Type = "room"
		event.Json = r
		data, _ := json.Marshal(event)
		fmt.Fprintln(conn, string(data))
	}
}

func admin(c model.Command, conn net.Conn, event model.Event) {
	switch c.Key {
	case "admin":
		changeUser(conn, "admin")
		SendCommand(event, conn, "admin_in")
		SendUser(conn, clients[conn].ID)
	case "out":
		changeUser(conn, "user")
		SendCommand(event, conn, "admin_out")
		SendUser(conn, clients[conn].ID)
	case "admin_err_key":
		SendCommand(event, conn, "admin_err_key")
	}
	return
}

func SendUser(conn net.Conn, id int) (m model.User) {
	var event model.Event
	usr, _ := json.Marshal(clients[conn])
	data := event.Prepare(usr, "user")
	fmt.Fprintln(conn, string(data))
	return
}

func SendAction(event model.Event, usr model.User, conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	data, _ := json.Marshal(event)
	for c, s := range clients {
		if c != conn && s.Room == usr.Room {
			fmt.Fprintln(c, string(data))
		}
	}
}

func changeUser(conn net.Conn, role string) {
	mutex.Lock()
	n := clients[conn]
	n.Role = role
	clients[conn] = n
	mutex.Unlock()
}

func Line(conn net.Conn, m model.Message) {
	usr := clients[conn]
	fmt.Printf("User <%s> with role <%s> send <%s> in room <%s>.\n", usr.Username, usr.Role, m.Msg, usr.Room)
}

func MakeId() (id int) {
	id = rand.Intn(50)
	for _, i := range clients {
		if i.ID == id {
			id = rand.Intn(len(clients) * 5)
		}
	}
	return
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
