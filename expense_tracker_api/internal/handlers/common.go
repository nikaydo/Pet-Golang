package handlers

import (
	"fmt"
	"log"
	env "main/internal/config"
	myjwt "main/internal/jwt"
	"main/internal/services"
	"net/http"
	"time"
)

type UserHandler struct {
	Service *services.DBServ
}

func NewUserHandler(service *services.DBServ) *UserHandler {
	return &UserHandler{Service: service}
}

func MakeCookie(name, value string, t time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(t),
		MaxAge:   int(t.Seconds()),
	}
}

func GetSubFromClaims(c *http.Cookie, env env.Config) (int, error) {
	j := myjwt.JwtTokens{AccessToken: c.Value, Env: env}
	if err := j.ValidateJwt(); err != nil {
		return 0, err
	}
	sub, ok := j.AccessClaims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("sub claim not found or invalid type")
	}
	return int(sub), nil
}

func handleServerError(w http.ResponseWriter, msg string, err error) {
	log.Println(msg+":", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
