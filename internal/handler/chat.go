package handler

import (
	"net/http"

	"github.com/RamzittoRamzotti/gotochat.git/internal/domain"
	"github.com/RamzittoRamzotti/gotochat.git/internal/services/chat"
	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	chatService *chat.ChatService
}

func NewChatHandler(chatService *chat.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (ch *ChatHandler) Chat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("name")
		if err != nil {
			http.Error(w, "missing name cookie", http.StatusBadRequest)
			return
		}

		roomID := chi.URLParam(r, "roomID")
		room, exists := ch.chatService.Rooms[roomID]
		if !exists {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := domain.NewClient(conn, cookie.Value)
		client.Rooms[roomID] = room
		room.Regsiter <- client

		go client.SendMessage(*ch.chatService.Logger)
		client.ReadMessage(*ch.chatService.Logger)
	}
}
