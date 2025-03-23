package command

import (
	"fmt"
	model "tcp/model"
)

func checkPermisson(role string) bool {
	if role != "admin" {
		fmt.Println("У вас нет прав администратора")
		return false
	}
	return true
}

func Room(list []string, s model.Send) model.Send {
	switch list[1] {

	case "add":
		s = addParam(list[2], s)
	case "del":
		s = delParam(list[2], s)
	case "list":
		listParam(s.Message.Room.Rooms)
	case "info":
		infoParam(s.Message.Room.CurrentRoom)
	case list[1]:
		for _, i := range s.Message.Room.Rooms {
			if i == list[1] {
				if i == s.Message.Room.CurrentRoom {
					fmt.Println("Вы уже находитесь в этой комнате")
				}
				s.Message.Room.CurrentRoom = i
			}
		}
	default:
		fmt.Println("Команда не распознана")
	}
	return s
}

func addParam(l string, s model.Send) model.Send {
	if !checkPermisson(s.Message.Role) {
		return s
	}
	fmt.Println("Комната", l, "добавлена")
	s.Message.Room.Rooms = append(s.Message.Room.Rooms, l)
	return s
}

func delParam(l string, s model.Send) model.Send {
	if !checkPermisson(s.Message.Role) {
		return s
	}
	for i, v := range s.Message.Room.Rooms {
		if v == l {
			s.Message.Room.Rooms = append(s.Message.Room.Rooms[:i], s.Message.Room.Rooms[i+1:]...)
		}
	}
	s.Message.Room.CurrentRoom = "main"
	fmt.Println("Комната", l, "удалена")
	return s
}

func listParam(s []string) {
	fmt.Println("Комнаты:")
	for _, i := range s {
		fmt.Println(i)
	}
}

func infoParam(s string) {
	fmt.Println("Комната", s)
}
