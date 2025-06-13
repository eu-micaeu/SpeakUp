package models

type Message struct {
	ID        string `json:"id" bson:"id"`
	ChatID    string `json:"chat_id" bson:"chat_id"`
	Content   string `json:"content" bson:"content"`
	Sender    string `json:"sender" bson:"sender"`
	Type      string `json:"type" bson:"type"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}
