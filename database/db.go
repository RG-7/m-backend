package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB(mongoURI string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("❌ MongoDB Connection Failed: ", err)
	}

	// Check if the database connection is alive
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB Ping Failed: ", err)
	}

	// Assign the client to the global variable
	Client = client
	log.Println("✅ Connected to MongoDB!")
}
