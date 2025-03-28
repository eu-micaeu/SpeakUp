package models

type Message struct {
	ID        string `json:"id"`
	ChatID    string `json:"chat_id"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	Type	  string `json:"type"`
	CreatedAt string `json:"timestamp"`
}
