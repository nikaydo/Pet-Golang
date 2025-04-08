package handlers

import (
	"log"
	myjwt "main/internal/jwt"
	"net/http"
	"strconv"
)

// DelTrans
// @Summary DelTrans
// @Description DelTrans
// @Tags user
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/deleteTrans  [delete]
func (h *UserHandler) DelTrans(w http.ResponseWriter, r *http.Request) {
	c, cfg, err := getFrom(r)
	if err != nil {
		handleServerError(w, err.Error(), err)
		return
	}
	var n []int
	id := r.URL.Query()["id"]
	for _, i := range id {
		i, err := strconv.Atoi(i)
		if err != nil {
			log.Println("Error converting string to int:", err)
		}
		n = append(n, i)
	}
	user_id, err := GetSubFromClaims(c, cfg)
	if err != nil {
		if err != myjwt.ErrTokenExpired {
			handleServerError(w, err.Error(), err)
		}
	}
	if err = h.Service.DelTrans(user_id, n); err != nil {
		handleServerError(w, "Error deleting transactions", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
