package voteendpoint

import (
	"context"
	"pollie"
	"pollie/middleware"
	"pollie/models"
	"pollie/pkg/vote/votesvc"
	"time"

	"github.com/go-kit/kit/endpoint"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Set struct {
	Vote endpoint.Endpoint
}

func NewSet(svc votesvc.Service) Set {
	voteEndpoint := makeEndpointVote(svc)
	voteEndpoint = middleware.SetUserID(voteEndpoint)

	return Set{
		Vote: voteEndpoint,
	}
}

func makeEndpointVote(svc votesvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Vote)

		// get the poll for the vote
		poll, err := svc.GetPoll(ctx, req.PollID)
		if err != nil {
			return pollie.Response{Err: err}, nil
		}
		// save the poll in the context
		ctx = context.WithValue(ctx, pollie.ContextKey("poll"), poll)

		var (
			voterID primitive.ObjectID
			userID  = middleware.GetUserID(ctx)
		)

		if userID != "" {
			voterID, _ = primitive.ObjectIDFromHex(userID)
		}

		v := models.Vote{
			Option:  req.Option,
			PollID:  poll.ID,
			VoterID: voterID,
			Meta: models.Meta{
				// set the initial meta values
				IP:        req.IP,
				Device:    req.Device,
				UserAgent: req.UserAgent,
			},
			CreatedAt: time.Now(),
		}

		p, err := svc.Vote(ctx, v)
		if err != nil {
			return pollie.Response{Err: err}, nil
		}

		return pollie.Response{
			Status: "success",
			Data: map[string]interface{}{
				"message": "vote successful",
				"poll":    p,
			},
		}, nil
	}
}
