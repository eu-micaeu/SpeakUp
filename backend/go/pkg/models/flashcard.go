package models

import "time"

type Flashcard struct {
	ID              string     `json:"id" bson:"_id,omitempty"`
	UserID          string     `json:"user_id" bson:"user_id"`
	Front           string     `json:"front" bson:"front"`
	Back            string     `json:"back" bson:"back"`
	ContextSentence string     `json:"context_sentence" bson:"context_sentence"`
	Explanation     string     `json:"explanation" bson:"explanation"`
	EaseFactor      float64    `json:"ease_factor" bson:"ease_factor"`
	Interval        int        `json:"interval" bson:"interval"`
	Repetitions     int        `json:"repetitions" bson:"repetitions"`
	NextReview      time.Time  `json:"next_review" bson:"next_review"`
	CreatedAt       time.Time  `json:"created_at" bson:"created_at"`
	LastReviewedAt  *time.Time `json:"last_reviewed_at,omitempty" bson:"last_reviewed_at,omitempty"`
}
