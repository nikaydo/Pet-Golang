package command

import (
	"fmt"
	model "tcp/model"
)

func Admin(list []string, m model.Message) string {
	switch list[1] {
	case "get":
		if len(list) < 3 {
			fmt.Println("Введите ключ прав администратора")
			return m.Role
		}
		if list[2] == "admin" {
			m.Role = "admin"
			fmt.Println("Права администратора получены")
			return m.Role
		}
		fmt.Println("Неверный ключь прав администратора")
	case "out":
		m.Role = "user"
	default:
		fmt.Println("Команда не распознана")
	}
	return m.Role
}
