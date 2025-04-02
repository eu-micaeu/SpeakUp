package models

type User struct {
	ID        string `bson:"id"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Email     string `bson:"email"`
	Password  string `bson:"password"`
	Language  string `bson:"language"`
}
