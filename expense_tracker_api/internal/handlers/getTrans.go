package handlers

import (
	"encoding/json"
	myjwt "main/internal/jwt"
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
