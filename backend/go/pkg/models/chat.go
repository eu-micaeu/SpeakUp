package models

type Chat struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	LegacyID  string `json:"-" bson:"id,omitempty"`
	UserID    string `json:"user_id" bson:"user_id"`
	StartTime string `json:"start_time" bson:"start_time"`
	Topic     string `json:"topic" bson:"topic"`
}

func (c *Chat) Normalize() {
	if c.LegacyID != "" {
		c.ID = c.LegacyID
	}
}
