package handlers

import (
	"encoding/json"
	myjwt "main/internal/jwt"
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
	c, cfg, err := getFrom(r)
	if err != nil {
		handleServerError(w, err.Error(), err)
		return
	}
	id, err := GetSubFromClaims(c, cfg)
	if err != nil {
		if err != myjwt.ErrTokenExpired {
			handleServerError(w, err.Error(), err)
		}
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
