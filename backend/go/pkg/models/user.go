package models

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	LegacyID string `json:"-" bson:"id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Language string `json:"language" bson:"language"`
	Level    string `json:"level" bson:"level"`
}

func (u *User) Normalize() {
	if u.ID == "" && u.LegacyID != "" {
		u.ID = u.LegacyID
	}
}
