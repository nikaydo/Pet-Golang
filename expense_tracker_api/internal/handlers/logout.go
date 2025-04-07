package handlers

import (
	"log"
	"net/http"
)

// Logout
// @Summary Log out
// @Description Invalidate the user's JWT token by setting an empty token in the cookie
// @Tags Auth
// @Success 200 {string} string "User signed out successfully"
// @Router /logout [post]

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("jwt")
	if err != nil {
		log.Println("Error getting cookie:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, MakeCookie("jwt", "", 0))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User signed out successfully"))
}
