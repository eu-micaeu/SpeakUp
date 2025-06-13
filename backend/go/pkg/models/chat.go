package models

type Chat struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"user_id" bson:"user_id"`
	StartTime string `json:"start_time" bson:"start_time"`
	Topic     string `json:"topic" bson:"topic"`
}
