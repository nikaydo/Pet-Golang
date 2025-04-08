package handlers

import (
	"encoding/json"
	"fmt"
	env "main/internal/config"
	"main/internal/models"
	"net/http"
)

// Register
// @Summary Register
// @Description Register and add in database
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body models.Auth true "Auth"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Router /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth
	var response string = "Username already exists"
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cfg, ok := r.Context().Value("config").(env.Config)
	if !ok {
		handleServerError(w, "Error checking user existence", fmt.Errorf("config missing in context"))
		return
	}
	exist, u, err := h.Service.IsUserExists(auth, cfg)
	if err != nil {
		handleServerError(w, "Error checking user existence", err)
		return
	}
	if !exist {
		if err = h.Service.AddUser(u, models.Balance{UserID: u.ID}, cfg); err != nil {
			handleServerError(w, "Failed to create user", err)
		}
		w.WriteHeader(http.StatusOK)
		response = "User created successfully"
	}
	w.WriteHeader(http.StatusConflict)
	w.Write([]byte(response))
}
