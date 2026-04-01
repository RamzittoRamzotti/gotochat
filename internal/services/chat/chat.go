package chat

import (
	"errors"
	"log/slog"

	"github.com/RamzittoRamzotti/gotochat.git/internal/domain"
	"github.com/RamzittoRamzotti/gotochat.git/internal/metrics"
)

var ErrRoomNotFound = errors.New("room not found")

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

func (s *ChatService) Room(roomID string) (*domain.Room, bool) {
	room, ok := s.Rooms[roomID]
	return room, ok
}

func (s *ChatService) JoinRoom(roomID string, client *domain.Client) error {
	room, ok := s.Rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}
	s.Clients[client.Name] = client
	client.Rooms[roomID] = room
	room.Regsiter <- client
	return nil
}

func (s *ChatService) LeaveClient(client *domain.Client) {
	delete(s.Clients, client.Name)
}
