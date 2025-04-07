package handlers

import (
	"encoding/json"
	env "main/internal/config"
	"net/http"
)

// GetTransactions
// @Summary GetTransactions
// @Description GetTransactions
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} models.Transaction
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/transactions [get]
func (h *UserHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
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
	tList, err := h.Service.Transactions(id)
	if err != nil {
		handleServerError(w, "Error getting transactions", err)
		return
	}
	jsonData, err := json.Marshal(tList)
	if err != nil {
		handleServerError(w, "Error encoding JSON", err)
		return
	}
	w.Write([]byte(jsonData))
}
