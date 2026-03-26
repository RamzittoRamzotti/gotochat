package app

import "github.com/RamzittoRamzotti/gotochat.git/internal/models/room"

type App struct {
	Rooms map[string]*room.Room
}
