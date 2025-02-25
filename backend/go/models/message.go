package models

type Message struct {
	ID        string    `json:"id"`
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
