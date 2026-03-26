package room

import (
	"github.com/RamzittoRamzotti/gotochat.git/internal/models/client"
	"github.com/RamzittoRamzotti/gotochat.git/internal/models/message"
)

type Room struct {
	ID       string
	Messages []*message.Message
	Clients  map[string]*client.Client
}
