package storage

import (
	"context"
	"errors"
	"log"
	"time"

	movement_record "github.com/WodBoard/models/movement_record/go"
	"gopkg.in/mgo.v2/bson"
)

// MovementRecordByEmail is a struct that contains email binding for a movement
// record to find it in storage
type MovementRecordByEmail struct {
	Email          string                         `bson:"email"`
	Name           string                         `bson:"name"`
	MovementRecord movement_record.MovementRecord `bson:"movement_record"`
}

// ListMovementRecords fetches a list of movements records by its email in mongo
func (s *Storage) ListMovementRecords(ctx context.Context, email string) ([]*movement_record.MovementRecord, error) {
	var res []*movement_record.MovementRecord

	ctx, _ = context.WithTimeout(ctx, time.Second*4)
	movementRecords := s.usersDatabase.Collection("movement_records")
	if movementRecords == nil {
		return nil, errors.New("error: couldn't read movement_records collection")
	}
	c, err := movementRecords.Find(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	defer c.Close(ctx)
	for c.Next(ctx) {
		var m MovementRecordByEmail
		err := c.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, &m.MovementRecord)
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return res, nil
}

// GetMovementRecord fetches one training by its name
func (s *Storage) GetMovementRecord(ctx context.Context, email, name string) (*movement_record.MovementRecord, error) {
	var res movement_record.MovementRecord

	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	movementRecords := s.usersDatabase.Collection("movement_records")
	if movementRecords == nil {
		return nil, errors.New("error: couldn't fetch movement records collection")
	}
	err := movementRecords.FindOne(ctx, bson.M{
		"email": email,
		"name":  name,
	}).
		Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// InsertMovementRecord inserts a new training in mongo bound by user email
func (s *Storage) InsertMovementRecord(ctx context.Context, email string, movementRecord movement_record.MovementRecord) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	movementRecords := s.usersDatabase.Collection("movement_records")
	if movementRecords == nil {
		return errors.New("error: couldn't insert movement_records collection")
	}
	_, err := movementRecords.InsertOne(ctx, MovementRecordByEmail{
		Email:          email,
		Name:           movementRecord.GetName(),
		MovementRecord: movementRecord,
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteMovementRecord deletes an existing movement record from mongo
func (s *Storage) DeleteMovementRecord(ctx context.Context, email string, name string) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	movementRecords := s.usersDatabase.Collection("movement_records")
	if movementRecords == nil {
		return errors.New("error: couldn't insert movement_records collection")
	}
	_, err := movementRecords.DeleteOne(ctx, bson.M{
		"email": email,
		"name":  name,
	})
	if err != nil {
		return err
	}
	return nil
}
