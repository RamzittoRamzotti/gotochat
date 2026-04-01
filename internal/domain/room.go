package domain

import (
	"github.com/RamzittoRamzotti/gotochat.git/internal/metrics"
)

type Room struct {
	ID         string
	Clients    map[string]*Client
	Regsiter   chan *Client
	Unregister chan *Client
	Messages   chan *Message
}

func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[string]*Client),
		Regsiter:   make(chan *Client, 10),
		Unregister: make(chan *Client, 10),
		Messages:   make(chan *Message, 10),
	}
}

func (r *Room) RemoveClient(client *Client) {
	delete(r.Clients, client.Name)
}

func (r *Room) Poll() {
	for {
		select {
		case client := <-r.Regsiter:
			r.Clients[client.Name] = client
			metrics.ActiveClients.Inc()
		case client := <-r.Unregister:
			metrics.ActiveClients.Dec()
			delete(r.Clients, client.Name)
		case message := <-r.Messages:
			for _, client := range r.Clients {
				if client.Name == message.SenderName {
					continue
				}
				go func() {
					client.SendChan <- message
				}()
			}
		}
	}
}
