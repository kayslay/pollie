package voterepo

import (
	"context"
	"fmt"
	"log"
	"pollie/config"
	"pollie/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	// Vote adds a vote to a poll
	Vote(v models.Vote) (interface{}, error)
	// IPExists checks if an ip has voted for a poll with the id
	// return a bool
	IPExists(id, ip string) (bool, error)
	// UpdatePollSummary
	UpdatePollSummary(eID string, vOption []int) error
}

// NewRepository create new repository
func NewRepository(mgo config.MgoFn) Repository {
	return repository{mgo: mgo}
}

type repository struct {
	mgo config.MgoFn
}

// IPExists check if a vote with ip exists
func (r repository) IPExists(eID, ip string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	c := r.mgo(models.VoteCollection)
	id, err := primitive.ObjectIDFromHex(eID)
	if err != nil {
		return false, models.ErrInvalidID
	}

	n, err := c.CountDocuments(ctx, bson.D{
		{"poll_id", id},
		{"meta.ip", ip},
	})

	if err != nil {
		return false, err
	}

	log.Println(n, ip)

	return n >= 1, nil
}

func (r repository) Vote(v models.Vote) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	c := r.mgo(models.VoteCollection)

	v.ID = primitive.NewObjectID()

	inserted, err := c.InsertOne(ctx, v)

	return inserted.InsertedID, err
}

// UpdatePollSummary update the summary of an poll
func (r repository) UpdatePollSummary(eID string, vOption []int) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	c := r.mgo(models.PollCollection)

	_eID, err := primitive.ObjectIDFromHex(eID)
	if err != nil {
		return models.ErrInvalidID
	}

	updateM := bson.M{
		"summary.votes": 1,
	}

	// pass the option_count index to be incremented
	for _, v := range vOption {
		updateM[fmt.Sprintf("summary.option_count.%d", v)] = 1
	}

	_, err = c.UpdateOne(ctx,
		bson.M{"_id": _eID},
		bson.M{
			"$set": bson.M{
				"summary.last_vote_at": time.Now(),
			},
			"$inc": updateM},
	)

	return err
}
