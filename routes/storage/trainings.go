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
	Name     string            `bson:"name"`
	Training training.Training `bson:"training"`
}

// ListTrainings fetches a list of trainings by its email in mongo
func (s *Storage) ListTrainings(ctx context.Context, email string) ([]*training.Training, error) {
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

// GetTraining fetches one training by its name
func (s *Storage) GetTraining(ctx context.Context, email, name string) (*training.Training, error) {
	var res training.Training

	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	users := s.usersDatabase.Collection("trainings")
	if users == nil {
		return nil, errors.New("error: couldn't fetch login collection")
	}
	err := users.FindOne(ctx, bson.M{
		"email": email,
		"name":  name,
	}).
		Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
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
		Name:     training.GetName(),
		Training: training,
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteTraining deletes an existing training from mongo
func (s *Storage) DeleteTraining(ctx context.Context, email string, name string) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	trainings := s.usersDatabase.Collection("trainings")
	if trainings == nil {
		return errors.New("error: couldn't insert trainings collection")
	}
	_, err := trainings.DeleteOne(ctx, bson.M{
		"email": email,
		"name":  name,
	})
	if err != nil {
		return err
	}
	return nil
}
