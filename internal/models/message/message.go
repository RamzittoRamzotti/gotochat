package message

type Message struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	ClientID string `json:"client_id"`
}
