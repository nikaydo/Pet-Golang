package middleware

import (
	"context"
	"log"
	env "main/internal/config"
	h "main/internal/handlers"
	myjwt "main/internal/jwt"
	"main/internal/services"
	"net/http"
)

type Middleware struct {
	Service *services.DBServ
}

func (f Middleware) CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt")
		if err != nil {
			log.Println("Error getting cookie:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		cfg, ok := r.Context().Value("config").(env.Config)
		if !ok {
			http.Error(w, "Error getting config", http.StatusInternalServerError)
			return
		}
		j := myjwt.JwtTokens{AccessToken: c.Value, Env: cfg}
		err = j.ValidateJwt()
		if err != nil {
			if err == myjwt.ErrTokenExpired {
				name, ok := j.AccessClaims["username"].(string)
				if !ok {
					log.Println("Error getting username:", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				u, err := f.Service.GetUser(name, cfg)
				if err != nil {
					log.Println("Error getting refresh token:", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				j.RefreshToken = u.Refresh
				err = j.ValidateRefresh()
				if err != nil {
					log.Println("Error validating refresh token:", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = j.CreateTokens(u.ID, u.Username, "user")
				if err != nil {
					log.Println("Error creating refresh token:", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = f.Service.UpdateRefreshToken(u, cfg)
				if err != nil {
					log.Println("Error updating refresh token:", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, h.MakeCookie("jwt", j.AccessToken, j.Env.TTL.Cookie))
				next.ServeHTTP(w, r)
				log.Println("Token refreshed successfully")
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ConfigMiddleware(e env.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "config", e)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
