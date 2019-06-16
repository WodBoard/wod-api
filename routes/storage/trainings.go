package storage

import (
	"context"
	"errors"
	"time"

	training "github.com/WodBoard/models/training/go"
	"gopkg.in/mgo.v2/bson"
)

// TrainingByEmail is a struct that contains email binding for a training
// to find it in storage
type TrainingByEmail struct {
	Email    string
	Training training.Training
}

// GetTrainingsByEmail fetches a list of trainings by its email in mongo
func (s *Storage) GetTrainingsByEmail(ctx context.Context, email string) ([]*training.Training, error) {
	var res []*training.Training
	var mongoRes []TrainingByEmail

	ctx, _ = context.WithTimeout(ctx, time.Second*4)
	trainings := s.usersDatabase.Collection("trainings")
	if trainings == nil {
		return nil, errors.New("error: couldn't fetch users collection")
	}
	c, err := trainings.Find(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	err = c.Decode(mongoRes)
	if err != nil {
		return nil, err
	}
	for _, t := range mongoRes {
		res = append(res, &t.Training)
	}
	return res, nil
}

// InsertTraining inserts a new training in mongo bound by user email
func (s *Storage) InsertTraining(ctx context.Context, email string, training training.Training) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	trainings := s.usersDatabase.Collection("trainings")
	if trainings == nil {
		return errors.New("error: couldn't access users collection")
	}
	_, err := trainings.InsertOne(ctx, &TrainingByEmail{
		Email:    email,
		Training: training,
	})
	if err != nil {
		return err
	}
	return nil
}
