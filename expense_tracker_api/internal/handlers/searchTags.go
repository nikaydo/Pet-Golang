package handlers

import (
	"encoding/json"
	env "main/internal/config"
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
	names := r.URL.Query()["tag"]
	if len(names) == 0 {
		http.Error(w, "name param required", http.StatusBadRequest)
		return
	}
	tags, err := h.Service.SearchTags(id, names...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonData))
}
