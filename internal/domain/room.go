package domain

type Room struct {
	ID         string
	Messages   []*Message
	Clients    map[string]*Client
	Regsiter   chan *Client
	Unregister chan *Client
}

func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Messages:   []*Message{},
		Clients:    make(map[string]*Client),
		Regsiter:   make(chan *Client, 10),
		Unregister: make(chan *Client, 10),
	}
}

func (r *Room) Poll() {
	for {
		select {
		case client := <-r.Regsiter:
			r.Clients[client.ID] = client
		case client := <-r.Unregister:
			delete(r.Clients, client.ID)
		}
	}
}
