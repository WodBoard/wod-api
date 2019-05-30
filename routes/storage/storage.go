package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Storage contains the different methods to store data into mongodb
// and wraps the connection and communication with mongo instance.
type Storage struct {
	DBClient       *mongo.Client
	usersDatabase  *mongo.Database
	loginsDatabase *mongo.Database
}

// NewStorage returns a new instance of Storage structure
func NewStorage(mongoClient *mongo.Client) *Storage {
	return &Storage{
		DBClient:      mongoClient,
		usersDatabase: mongoClient.Database("users"),
	}
}
