package handlers

import (
	"encoding/json"
	myjwt "main/internal/jwt"
	"net/http"
)

// SearchTags godoc
// @Summary SearchTags
// @Description SearchTags
// @Tags user
// @Accept json
// @Param name query string true "name"
// @Produce json
// @Success 200
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/tag [get]
func (h *UserHandler) SearchTags(w http.ResponseWriter, r *http.Request) {
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
	tags := r.URL.Query()["tag"]
	if len(tags) == 0 {
		http.Error(w, "name param required", http.StatusBadRequest)
		return
	}
	trans, err := h.Service.SearchTags(id, tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonData))
}
