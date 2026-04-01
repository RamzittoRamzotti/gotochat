package chat

import (
	"log/slog"

	"github.com/RamzittoRamzotti/gotochat.git/internal/domain"
	"github.com/RamzittoRamzotti/gotochat.git/internal/metrics"
)

type ChatService struct {
	Rooms   map[string]*domain.Room
	Clients map[string]*domain.Client
	Logger  *slog.Logger
}

func NewChatService(logger *slog.Logger) *ChatService {
	return &ChatService{
		Rooms:   make(map[string]*domain.Room),
		Clients: make(map[string]*domain.Client),
		Logger:  logger,
	}
}

func (s *ChatService) AddRoom(room *domain.Room) {
	s.Rooms[room.ID] = room
	metrics.ActiveRooms.Inc()
}

func (s *ChatService) AddClientToRoom(roomID string, client *domain.Client) (*domain.Room, bool) {
	room, exists := s.Rooms[roomID]
	if !exists {
		return nil, false
	}
	s.Clients[client.Name] = client
	return room, true
}

func (s *ChatService) RemoveClientFromRoom(roomID string, client *domain.Client) {
	room, exists := s.Rooms[roomID]
	if !exists {
		return
	}
	delete(s.Clients, client.Name)
	room.RemoveClient(client)
}
