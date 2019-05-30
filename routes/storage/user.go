package storage

import (
	"context"
	"errors"
	"time"

	user "github.com/WodBoard/models/user/go"
	"gopkg.in/mgo.v2/bson"
)

// GetUserByEmail fetches a user by its email in mongo
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var res user.User

	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	users := s.usersDatabase.Collection("users")
	if users == nil {
		return nil, errors.New("error: couldn't fetch users collection")
	}
	err := users.FindOne(ctx, bson.M{"username": email}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// InsertUser inserts a new user in mongo
func (s *Storage) InsertUser(ctx context.Context, user *user.User) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	users := s.usersDatabase.Collection("users")
	if users == nil {
		return errors.New("error: couldn't access users collection")
	}
	_, err := users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
