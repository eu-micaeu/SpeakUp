package models

import "time"

type User struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	LegacyID       string    `json:"-" bson:"id,omitempty"`
	Name           string    `json:"name" bson:"name"`
	Email          string    `json:"email" bson:"email"`
	Password       string    `json:"password" bson:"password"`
	Language       string    `json:"language" bson:"language"`
	Level          string    `json:"level" bson:"level"`
	ResetCodeHash  string    `json:"-" bson:"reset_code_hash,omitempty"`
	ResetExpiresAt time.Time `json:"-" bson:"reset_expires_at,omitempty"`
}

func (u *User) Normalize() {
	if u.LegacyID != "" {
		u.ID = u.LegacyID
	}
}

