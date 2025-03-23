package command

import (
	"fmt"
	"net"
	"strings"

	model "tcp/model"
)

var (
	ErrNoParam = fmt.Errorf("нет параметров")
)

func Command(s model.Send, conn net.Conn) (model.Send, error) {
	msg := strings.Replace(s.Message.Msg, "!", "", 1)
	list := strings.Split(msg, " ")
	switch list[0] {
	case "help":
		cHelp()
	case "info":
		cInfo(s.Message, conn)
	case "room":
		if len(list) < 2 {
			return s, ErrNoParam
		}
		n := Room(list, s)
		s.Message.Room = n.Message.Room
	case "admin":
		if len(list) < 2 {
			return s, ErrNoParam
		}
		s.Message.Role = Admin(list, s.Message)
	default:
		fmt.Println("Команда не распознана")
	}
	return s, nil
}

func cHelp() {
	fmt.Println("help - список команд")
	fmt.Println("info - информация о пользователе")
	fmt.Println("room - управление комнатами\n    [название комнаты] - переход в комнату\n    add  - добавить комнату\n    del  - удалить комнату\n    list - список комнат\n    info - информация о комнате")
	fmt.Println("admin - режим администратора\n    get [key]  - получить права администратора\n    out  - выйти из режима администратора")
}

func cInfo(m model.Message, conn net.Conn) {
	fmt.Println("Локальный IP", conn.LocalAddr().String())
	fmt.Println("Удалённый IP", conn.RemoteAddr().String())
	fmt.Println("Пользователь", m.Username)
	fmt.Println("Комната", m.Room)
}
