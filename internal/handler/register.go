package handler

import (
	"encoding/json"
	"net/http"

	models "github.com/RamzittoRamzotti/gotochat.git/internal/domain/models"
	auth "github.com/RamzittoRamzotti/gotochat.git/internal/services/auth"
)

type RegisterHandler struct {
	authService *auth.AuthService
}

func NewRegisterHandler(authService *auth.AuthService) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
	}
}

func (h *RegisterHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		client := &models.UserCreate{}

		err := json.NewDecoder(r.Body).Decode(&client)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = h.authService.Register(client.Username, client.Password)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
