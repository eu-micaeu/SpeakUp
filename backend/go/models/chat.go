package models

type Chat struct {
	ID        string    `json:"id"`
	UserID    string `json:"user_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Topic     string `json:"topic"`
}
