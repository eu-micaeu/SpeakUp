package models

type Word struct {
	ID              string  `json:"id" bson:"id"`
	UserID          string  `json:"user_id" bson:"user_id"`
	Word            string  `json:"word" bson:"word"`
	WordTranslated  string  `json:"word_translated" bson:"word_translated"`
	ContextInPhrase string  `json:"context_in_phrase" bson:"context_in_phrase"`
	DeckID          *string `json:"deck_id,omitempty" bson:"deck_id,omitempty"`
	Language        string  `json:"language" bson:"language"`
	Level           string  `json:"level" bson:"level"`
	CreatedAt       string  `json:"created_at" bson:"created_at"`
}