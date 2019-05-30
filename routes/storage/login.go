package storage

import (
	"context"
	"errors"
	"time"

	user "github.com/WodBoard/models/user/go"
	"gopkg.in/mgo.v2/bson"
)

// GetLoginByEmail fetches login informations by email in mongo
func (s *Storage) GetLoginByEmail(ctx context.Context, email string) (*user.Login, error) {
	var res user.Login

	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	users := s.usersDatabase.Collection("logins")
	if users == nil {
		return nil, errors.New("error: couldn't fetch login collection")
	}
	err := users.FindOne(ctx, bson.M{"username": email}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// InsertLogin inserts a new login informations about a user in mongo
func (s *Storage) InsertLogin(ctx context.Context, login *user.Login) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	users := s.usersDatabase.Collection("logins")
	if users == nil {
		return errors.New("error: couldn't access login collection")
	}
	_, err := users.InsertOne(ctx, login)
	if err != nil {
		return err
	}
	return nil
}
