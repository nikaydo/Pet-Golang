package router

import (
	env "main/internal/config"
	handler "main/internal/handlers"
	m "main/internal/middleware"
	rep "main/internal/repositories"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Router(e env.Config, handler *handler.UserHandler, userRepo *rep.File) http.Handler {

	mid := m.Middleware{File: userRepo}
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
	})
	r.Post("/auth", handler.AuthUser)
	return r
}
