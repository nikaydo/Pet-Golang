package handlers

import (
	"encoding/json"
	myjwt "main/internal/jwt"
	"main/internal/models"
	"net/http"
)

// MakeTransactions
// @Summary MakeTransactions
// @Description MakeTransactions
// @Tags user
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction"
// @Success 201 {object} models.Transaction
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/newtransaction [post]
func (h *UserHandler) MakeTransactions(w http.ResponseWriter, r *http.Request) {
	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		handleServerError(w, "Invalid request body", err)
		return
	}
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

	t.UserID = id
	if err = h.Service.NewTransactions(t); err != nil {
		handleServerError(w, "Failed to create transaction", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
