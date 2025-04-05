package main

import (
	"fmt"
	env "main/internal/config"
	db "main/internal/db"
	handler "main/internal/handlers"
	rep "main/internal/repositories"
	rt "main/internal/router"
	services "main/internal/services"
	"net/http"
)

func main() {
	e := env.GetConfig()
	database := db.InitDB(e.DB.Path)
	userRepo := rep.NewRepository(database)
	err := userRepo.MakeTable()
	if err != nil {
		panic(err)
	}
	userService := services.NewUserService(userRepo)
	h := handler.NewUserHandler(userService)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", e.Server.Host, e.Server.Port), rt.Router(e, h, userRepo)); err != nil {
		fmt.Printf("Пу пу пу: %s", err.Error())
		return
	}
}
