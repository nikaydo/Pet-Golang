package handlers

import (
	"encoding/json"
	env "main/internal/config"
	"net/http"
)

// GetBalance
// @Summary GetBalance
// @Description GetBalance
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.Balance
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/balance [get]
func (h *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwt")
	if err != nil {
		handleServerError(w, "Unauthorized: token not found", err)
		return
	}
	cfg, ok := r.Context().Value("config").(env.Config)
	if !ok {
		handleServerError(w, "could not get config", err)
		return
	}
	id, err := GetSubFromClaims(c, cfg)
	if err != nil {
		handleServerError(w, err.Error(), err)
		return
	}
	money, err := h.Service.Balance(id)
	if err != nil {
		handleServerError(w, err.Error(), err)
		return
	}
	jsonData, err := json.Marshal(money)
	if err != nil {
		handleServerError(w, "Error marshalling response", err)
		return
	}
	w.Write([]byte(jsonData))
}
