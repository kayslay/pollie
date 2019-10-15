package pollrepo

import (
	"context"
	"errors"
	"pollie"
	"pollie/config"
	"pollie/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	Create(ctx context.Context, p models.Poll) error
	Delete(ctx context.Context, id, uID string) error
	GetMany(ctx context.Context, uID string) ([]models.Poll, error)
	Get(ctx context.Context, id string) (models.Poll, error)
}

// NewRepository create new repository
func NewRepository(mgo config.MgoFn) Repository {
	return repository{mgo: mgo}
}

type repository struct {
	mgo config.MgoFn
}

func (r repository) Create(_ctx context.Context, p models.Poll) error {
	ctx, _ := context.WithTimeout(_ctx, 5*time.Second)
	c := r.mgo(models.PollCollection)
	p.ID = primitive.NewObjectID()

	_, err := c.InsertOne(ctx, p)

	return pollie.WrapErr(err, "create_poll")
}

func (r repository) Delete(_ctx context.Context, id, uID string) error {
	ctx, _ := context.WithTimeout(_ctx, 5*time.Second)
	c := r.mgo(models.PollCollection)
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidID
	}
	_uID, err := primitive.ObjectIDFromHex(uID)
	if err != nil {
		return models.ErrInvalidID
	}

	result, err := c.UpdateOne(ctx, bson.M{
		"_id":        _id,
		"user_id":    _uID,
		"deleted_at": nil,
	}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})

	// the result matched was not equal to 1. means nothing was updated
	if result.MatchedCount != 1 {
		return errors.New("could not find election to to delete")
	}

	return pollie.WrapErr(err, "delete_poll")
}

func (r repository) GetMany(_ctx context.Context, uID string) ([]models.Poll, error) {

	ctx, _ := context.WithTimeout(_ctx, 5*time.Second)
	c := r.mgo(models.PollCollection)
	_uID, err := primitive.ObjectIDFromHex(uID)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	pp := []models.Poll{}

	cur, err := c.Find(ctx, bson.M{"user_id": _uID, "deleted_at": nil, "expired_at": nil})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var p models.Poll
		err := cur.Decode(&p)
		if err != nil {
			return nil, pollie.WrapErr(err, "poll_getmany_cursor")
		}
		pp = append(pp, p)
	}

	return pp, nil
}

func (r repository) Get(_ctx context.Context, id string) (models.Poll, error) {
	ctx, _ := context.WithTimeout(_ctx, 5*time.Second)
	c := r.mgo(models.PollCollection)
	p := models.Poll{}
	// the default filter uses the short code to get the poll
	filter := bson.M{"short_code": id}

	_id, err := primitive.ObjectIDFromHex(id)
	// if the id is a valid ObjectID user the _id to the poll
	if err == nil {
		filter = bson.M{"_id": _id}
	}

	sResult := c.FindOne(ctx, filter)
	if err := sResult.Err(); err != nil {

		if err == mongo.ErrNoDocuments {
			return p, errors.New("poll does not exist")
		}
		return p, err
	}

	err = sResult.Decode(&p)
	if err != nil {
		return p, pollie.WrapErr(err, "poll_get_cursor")
	}

	return p, err
}
