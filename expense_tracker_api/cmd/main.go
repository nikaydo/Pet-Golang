package main

import (
	"fmt"
	env "main/internal/config"
	handler "main/internal/handlers"
	rep "main/internal/repositories"
	rt "main/internal/router"
	services "main/internal/services"
	"net/http"
)

// @title           Expense Tracker API
// @version         1.0
// @description     API for expense tracker
// @host      localhost:8080
func main() {
	e := env.GetConfig()
	Database := rep.NewRepository(&e)
	err := Database.MakeTable()
	if err != nil {
		panic(err)
	}
	Serv := services.NewDBServ(Database)
	h := handler.NewUserHandler(Serv)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", e.Server.Host, e.Server.Port), rt.Router(e, h, Serv)); err != nil {
		fmt.Printf("Пу пу пу: %s", err.Error())
		return
	}
}
