package models

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"timestamp"`
}
