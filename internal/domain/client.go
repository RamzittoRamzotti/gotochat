package domain

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/RamzittoRamzotti/gotochat.git/internal/metrics"
	"github.com/gorilla/websocket"
)

type Client struct {
	Name     string
	Conn     *websocket.Conn
	Rooms    map[string]*Room
	GetChan  chan *Message
	SendChan chan *Message
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func NewClient(conn *websocket.Conn, name string) *Client {
	return &Client{
		Conn:     conn,
		Rooms:    make(map[string]*Room),
		SendChan: make(chan *Message, 10),
		Name:     name,
	}
}

func (c *Client) SendMessage(logger slog.Logger) error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.SendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}
			mes, err := json.Marshal(message)
			if err != nil {
				logger.Error("error while marshaling message", "error", err)
				return err
			}
			w.Write(mes)
			metrics.MessagesTotal.Inc()

			if err := w.Close(); err != nil {
				return err
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}

func (c *Client) ReadMessage(logger slog.Logger) error {
	defer func() {
		metrics.ActiveClients.Dec()
		for _, room := range c.Rooms {
			room.Unregister <- c
		}
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("error while reading message", "error", err)
			}
			return err
		}
		for _, room := range c.Rooms {
			logger.Info("message received", "message", string(message), "room_id", room.ID, "client_name", c.Name)
			room.Messages <- &Message{
				ID:         c.Name + time.Now().String(),
				Content:    string(message),
				RoomID:     room.ID,
				SenderName: c.Name,
			}
		}
	}
}
