package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	VoteCollection = "votes"
)

type Vote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	PollID    primitive.ObjectID `json:"poll_id" bson:"poll_id"`
	VoterID   primitive.ObjectID `json:"voter_id" bson:"voter_id"`
	Option    []int              `json:"option" bson:"option"`
	Meta      Meta               `json:"meta" bson:"meta"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// 127.0.0.1:59942 127.0.0.1:59942
type Meta struct {
	IP          string  `json:"ip" bson:"ip"`
	CountryName string  `json:"country" bson:"country"`
	CountryCode string  `json:"country_code" bson:"country_code"`
	RegionName  string  `json:"state" bson:"state"`
	RegionCode  string  `json:"region_code" bson:"region_code"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	UserAgent   string  `json:"user_agent" bson:"user_agent"`
	Device      string  `json:"device" bson:"device"`
}
