package domain

type Message struct {
	ID         string `json:"client_id"`
	Content    string `json:"content"`
	RoomID     string `json:"room_id"`
	SenderName string `json:"sender_name"`
}
