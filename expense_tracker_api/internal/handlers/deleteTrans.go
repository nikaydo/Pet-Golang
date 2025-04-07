package handlers

import (
	env "main/internal/config"
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
	id := r.URL.Query()["id"]
	n, err := strconv.Atoi(id[0])
	if err != nil {
		handleServerError(w, err.Error(), err)
		return
	}
	user_id, err := GetSubFromClaims(c, cfg)
	if err != nil {
		handleServerError(w, err.Error(), err)
		return
	}
	if err = h.Service.DelTrans(user_id, n); err != nil {
		handleServerError(w, "Error deleting transactions", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
