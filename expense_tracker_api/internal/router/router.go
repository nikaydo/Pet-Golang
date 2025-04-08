package router

import (
	_ "main/docs"
	env "main/internal/config"
	handler "main/internal/handlers"
	m "main/internal/middleware"
	"main/internal/services"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(e env.Config, handler *handler.UserHandler, userRepo *services.DBServ) http.Handler {
	mid := m.Middleware{Service: userRepo}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(m.ConfigMiddleware(e))
	r.Route("/user", func(r chi.Router) {
		r.Use(mid.CheckJWT)
		r.Get("/balance", handler.GetBalance)               
		r.Get("/transactions", handler.GetTransactions)     
		r.Post("/newtransaction", handler.MakeTransactions) 
		r.Delete("/deleteTrans", handler.DelTrans)          
		r.Get("/logout", handler.Logout)                    
		r.Post("/tag", handler.SearchTags)                  
	})
	r.Post("/signin", handler.SignIn)            
	r.Post("/register", handler.Register)        
	r.Get("/swagger/*", httpSwagger.WrapHandler) 
	return r
}
