package models

type User struct {
    ID        string `json:"id" bson:"id"`
    Name      string `json:"name" bson:"name"`
    Email     string `json:"email" bson:"email"`
    Password  string `json:"password" bson:"password"`
    Language  string `json:"language" bson:"language"`
    Level     string `json:"level" bson:"level"`
}