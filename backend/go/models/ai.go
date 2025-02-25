package models

type AI struct {
	ID   string `json:"id" bson:"id"`
	ModelName string `json:"model_name" bson:"model_name"`
	Version string `json:"version" bson:"version"`
}