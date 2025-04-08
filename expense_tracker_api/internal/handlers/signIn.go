package handlers

import (
	"encoding/json"
	"fmt"
	env "main/internal/config"
	myjwt "main/internal/jwt"
	"main/internal/models"
	"net/http"
)

// SignIn
// @Summary Sign in
// @Description Authenticate a user and generate JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body models.Auth true "Auth"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Router /signin [post]
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth
	var response string = "User signed in successfully"
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
	exist, u, err := h.Service.IsUserExists(auth, cfg)
	if err != nil {
		handleServerError(w, "Error checking user existence", err)
		return
	}
	if !exist {
		w.WriteHeader(http.StatusConflict)
		response = "Username does not exist"
	}
	if err = j.CreateTokens(u.ID, u.Username, "user"); err != nil {
		handleServerError(w, "Failed to generate JWT tokens", err)
		return
	}
	u.Refresh = j.RefreshToken
	if err = h.Service.UpdateRefreshToken(u, cfg); err != nil {
		handleServerError(w, "Failed to update refresh token", err)
		return
	}
	http.SetCookie(w, MakeCookie("jwt", j.AccessToken, cfg.TTL.Cookie))
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(response))
}
