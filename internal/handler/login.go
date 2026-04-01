package handler

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/RamzittoRamzotti/gotochat.git/internal/services/auth"
	models "github.com/RamzittoRamzotti/gotochat.git/internal/domain/models"
)

type LoginHandler struct {
	authService *auth.AuthService
}

func NewLoginHandler(authService *auth.AuthService) *LoginHandler {
	return &LoginHandler{authService: authService}
}

func (h *LoginHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UserCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		token, err := h.authService.Login(req.Username, req.Password)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false, // браузер должен читать для WebSocket
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "name",
			Value:   req.Username,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
	}
}
