package main

import (
	"context"
	"log"
	"os"
	"recipes-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	users := map[string]models.User{
		"admin": {
			Username: "Admin",
			Email:    "admin@example.com",
			Password: "fCRmh4Q2J7Rseqkz",
		},
		"packt": {
			Username: "Packt",
			Email:    "packt@example.com",
			Password: "RE4zfHB35VPtTkbT",
		},
		"mlabouardy": {
			Username: "mlabouardy",
			Email:    "mlabouardy@example.com",
			Password: "L3nSFRcZzNQ67bcc",
		},
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	for _, user := range users {
		// user.Password.Plaintext =
		switch user.Username {
		case "Admin":
			hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			user.Password = string(hash)
		case "Packt":
			hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			user.Password = string(hash)
		case "mlabouardy":
			hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			user.Password = string(hash)
		}
		collection.InsertOne(ctx, bson.M{
			"username": user.Username,
			"password": user.Password,
			"email":    user.Email,
		})
	}
}
