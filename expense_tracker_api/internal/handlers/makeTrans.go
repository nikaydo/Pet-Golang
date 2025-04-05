package handlers

import (
	"encoding/json"
	"fmt"
	env "main/internal/config"
	"main/internal/models"
	"net/http"
)

func (h *UserHandler) MakeTransactions(w http.ResponseWriter, r *http.Request) {
	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		handleServerError(w, "Invalid request body", err)
		return
	}
	c, err := r.Cookie("jwt")
	if err != nil {
		handleServerError(w, "Unauthorized: token not found", err)
		return
	}
	cfg, ok := r.Context().Value("config").(env.Config)
	if !ok {
		handleServerError(w, "Could not get config", fmt.Errorf("config missing in context"))
		return
	}
	id, err := GetSubFromClaims(c, cfg)
	if err != nil {
		handleServerError(w, "Invalid token", err)
		return
	}
	t.UserID = id
	if err = h.Service.NewTransactions(t); err != nil {
		handleServerError(w, "Failed to create transaction", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
