package models

type User struct {
	Password string `json:"password" bson:"password"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}
