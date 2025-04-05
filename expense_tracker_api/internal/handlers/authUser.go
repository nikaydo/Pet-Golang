package handlers

import (
	"encoding/json"
	"fmt"
	env "main/internal/config"
	myjwt "main/internal/jwt"
	"main/internal/models"
	"net/http"
)

func (h *UserHandler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cfg, ok := r.Context().Value("config").(env.Config)
	if !ok {
		handleServerError(w, "Error checking user existence", fmt.Errorf("config missing in context"))
		return
	}
	j := myjwt.JwtTokens{Env: cfg}
	exist, u, err := h.Service.IsUserExists(auth)
	if err != nil {
		handleServerError(w, "Error checking user existence", err)
		return
	}
	if !exist {
		if err = h.Service.AddUser(u, models.Balance{UserID: u.ID}); err != nil {
			handleServerError(w, "Failed to create user", err)
		}
		_, u, err = h.Service.IsUserExists(auth)
		if err != nil {
			handleServerError(w, "Error checking user existence", err)
			return
		}
	}
	if err = j.CreateTokens(u.ID, u.Username, "user"); err != nil {
		handleServerError(w, "Failed to generate JWT tokens", err)
		return
	}
	if err = h.Service.UpdateRefreshToken(u); err != nil {
		handleServerError(w, "Failed to update refresh token", err)
		return
	}
	http.SetCookie(w, MakeCookie("jwt", j.AccessToken, cfg.TTL.AccessToken))
	w.WriteHeader(getStatusCode(exist))
	w.Write([]byte("Authentication successful"))
}
