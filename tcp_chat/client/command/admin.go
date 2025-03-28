package command

import "fmt"

func Admin(list []string) (key string) {

	switch list[1] {
	case "get":
		if len(list) < 3 {
			fmt.Println("Введите ключ прав администратора")
			return
		}
		if list[2] == "admin" {
			key = list[2]
		}
	case "out":
		key = "out"
	}
	return
}
