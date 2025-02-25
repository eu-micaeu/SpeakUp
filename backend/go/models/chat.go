package models

type Chat struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DifficultyLevel string `json:"difficulty_level"`
	Topic           string `json:"topic"`
}
