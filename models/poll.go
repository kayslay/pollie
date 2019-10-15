package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PollCollection poll collection name
const PollCollection = "polls"

const (
	// vote types

	PTypeSingle = "SINGLE"
)

const (
	// poll identity

	IDPayment       = "PAYMENT"
	IDUser          = "USER"
	IDIPAddr        = "IP_ADDRESS"
	IDCookie        = "COOKIE"
	IDSocialTwitter = "SOCIAL_TWITTER"
)

type Poll struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	UserID  primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type    string             `json:"type" bson:"type"`
	NilVote bool               `json:"nil_vote" bson:"nil_vote"`
	// NOTE Auth & Invisible does not seem to have a usecase at the moment
	Invisible   bool      `json:"invisible" bson:"invisible"`
	Auth        bool      `json:"auth" bson:"auth"`
	Description string    `json:"description" bson:"description"`
	Option      []string  `json:"option" bson:"option"`
	Tags        []string  `json:"tags" bson:"tags"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	ShortCode   string    `json:"short_code" bson:"short_code"`
	// Identity := user|payment|ip addresses| browser cookies| social (twitter/facebook)
	Identity string      `json:"identity" bson:"identity"`
	Summary  PollSummary `json:"summary" bson:"summary"`
	// using a *time.Time to pass in null values
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
}

type PollSummary struct {
	Votes       int64     `json:"votes" bson:"votes"`
	LastVoteAt  time.Time `json:"last_vote_at" bson:"last_vote_at"`
	OptionCount []int64   `json:"option_count" bson:"option_count"`
}
