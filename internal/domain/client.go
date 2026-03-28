package domain

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Rooms    map[string]*Room
	GetChan  chan *Message
	SendChan chan *Message
	ID       string
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func NewClient(conn *websocket.Conn, id string) *Client {
	return &Client{
		Conn:     conn,
		Rooms:    make(map[string]*Room),
		GetChan:  make(chan *Message, 10),
		SendChan: make(chan *Message, 10),
		ID:       id,
	}
}

func (c *Client) SendMessage(message *Message) error {
	defer func() {
		for _, room := range c.Rooms {
			room.UnregisterClient(c)
		}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		select {
		case <-time.After(pingPeriod):
		}

	}
}

func (c *Client) ReadMessage() (*Message, error) {
	var message Message
	err := c.Conn.ReadJSON(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
