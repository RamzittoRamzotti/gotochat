package app

import (
	"github.com/RamzittoRamzotti/gotochat.git/internal/models/room"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() *App {
	return &App{
		rooms: make(map[string]*room.Room),
	}
}
