package main

import (
	"context"
	"fmt"
	"os"

	"github.com/WodBoard/wod-api/routes"
	"github.com/WodBoard/wod-api/routes/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	mongoClient, err := mongo.NewClient(options.Client().
		ApplyURI(os.Getenv("MONGO_URI")).
		SetAuth(options.Credential{
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create mongodb instance", err)
		os.Exit(2)
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't connect to mongodb instance", err)
		os.Exit(2)
	}

	storage := storage.NewStorage(mongoClient)
	handler := routes.NewHandler(storage, os.Getenv("LISTEN_ADDR"))

	handler.HandleRoutes()
}
