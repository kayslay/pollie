package voteendpoint

import (
	"context"
	"pollie"
	"pollie/middleware"
	"pollie/models"
	"pollie/pkg/vote/votesvc"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Set struct {
	Vote    endpoint.Endpoint
	GetPoll endpoint.Endpoint
}
type pollResponse struct {
	models.Poll
	Summary string `json:"summary,omitempty" bson:"summary"`
}

func NewSet(svc votesvc.Service, logger log.Logger) Set {
	voteEndpoint := makeEndpointVote(svc)
	voteEndpoint = middleware.SetUserID(voteEndpoint)
	voteEndpoint = middleware.Logger(log.With(logger, "method", "make_vote"))(voteEndpoint)

	getPollEndpoint := makeEndpointGetPoll(svc)
	getPollEndpoint = middleware.SetUserID(getPollEndpoint)
	getPollEndpoint = middleware.Logger(log.With(logger, "method", "get_vote_poll"))(getPollEndpoint)

	return Set{
		Vote:    voteEndpoint,
		GetPoll: getPollEndpoint,
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

func makeEndpointGetPoll(svc votesvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)
		p, err := svc.GetPoll(ctx, req.ID)
		if err != nil {
			return pollie.Response{Err: err}, nil
		}

		return pollie.Response{
			Status: "success",
			Data: map[string]interface{}{
				"message": "successfully fetched poll",
				"poll":    pollResponse{Poll: p},
			},
		}, nil
	}
}
