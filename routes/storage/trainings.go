package storage

import (
	"context"
	"errors"
	"log"
	"time"

	training "github.com/WodBoard/models/training/go"
	"gopkg.in/mgo.v2/bson"
)

// TrainingByEmail is a struct that contains email binding for a training
// to find it in storage
type TrainingByEmail struct {
	Email    string            `bson:"email"`
	Training training.Training `bson:"training"`
}

// GetTrainingsByEmail fetches a list of trainings by its email in mongo
func (s *Storage) GetTrainingsByEmail(ctx context.Context, email string) ([]*training.Training, error) {
	var res []*training.Training

	ctx, _ = context.WithTimeout(ctx, time.Second*4)
	trainings := s.usersDatabase.Collection("trainings")
	if trainings == nil {
		return nil, errors.New("error: couldn't read trainings collection")
	}
	c, err := trainings.Find(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	defer c.Close(ctx)
	for c.Next(ctx) {
		var t TrainingByEmail
		err := c.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, &t.Training)
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return res, nil
}

// InsertTraining inserts a new training in mongo bound by user email
func (s *Storage) InsertTraining(ctx context.Context, email string, training training.Training) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	trainings := s.usersDatabase.Collection("trainings")
	if trainings == nil {
		return errors.New("error: couldn't insert trainings collection")
	}
	_, err := trainings.InsertOne(ctx, TrainingByEmail{
		Email:    email,
		Training: training,
	})
	if err != nil {
		return err
	}
	return nil
}
