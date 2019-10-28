package config

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type MgoFn func(string) *mongo.Collection

// NewMgoDB returns a func that returns a collection
func NewMgoDB() (MgoFn, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MGO_URL")))
	if err != nil {
		return nil, err
	}

	// set db
	db := client.Database(viper.GetString("MGO_DB"))

	go createIndex(ctx, db)
	// return mgo collection func
	return func(collName string) *mongo.Collection {
		return db.Collection(collName)
	}, nil
}

func createIndex(ctx context.Context, d *mongo.Database) {

	uniqueTrue := true
	// Poll index
	// shortcode
	pollShortCodeIndex := mongo.IndexModel{
		Keys: bson.M{
			"short_code": 1,
		},
		Options: &options.IndexOptions{
			Unique: &uniqueTrue,
		},
	}

	pollUserIDIndex := mongo.IndexModel{
		Keys: bson.M{
			"user_id": 1,
		},
	}

	pollExpiresAtIndex := mongo.IndexModel{
		Keys: bson.M{
			"expires_at": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(60),
	}

	d.Collection("polls").Indexes().CreateOne(ctx, pollShortCodeIndex)
	d.Collection("polls").Indexes().CreateOne(ctx, pollUserIDIndex)
	d.Collection("polls").Indexes().CreateOne(ctx, pollExpiresAtIndex)

	// votes
	voteUserIDIndex := mongo.IndexModel{
		Keys: bson.M{
			"user_id": 1,
		},
	}
	votePollIDIndex := mongo.IndexModel{
		Keys: bson.M{
			"poll_id": 1,
		},
	}
	voteIPIndex := mongo.IndexModel{
		Keys: bson.M{
			"poll_id": -1,
			"meta.ip": 1,
		},
	}

	d.Collection("votes").Indexes().CreateOne(ctx, voteUserIDIndex)
	d.Collection("votes").Indexes().CreateOne(ctx, votePollIDIndex)
	d.Collection("votes").Indexes().CreateOne(ctx, voteIPIndex)

	// ip
	ipIndex := mongo.IndexModel{
		Keys: bson.M{
			"ip": 1,
		},
	}

	d.Collection("ips").Indexes().CreateOne(ctx, ipIndex)

}
