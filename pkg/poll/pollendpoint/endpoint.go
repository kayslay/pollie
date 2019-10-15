package pollendpoint

import (
	"context"
	"errors"
	"fmt"
	"pollie"
	"pollie/models"

	"pollie/middleware"
	"pollie/pkg/poll/pollsvc"

	"github.com/go-kit/kit/endpoint"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Set struct {
	Create  endpoint.Endpoint
	Delete  endpoint.Endpoint
	GetMany endpoint.Endpoint
	Get     endpoint.Endpoint
}

func NewSet(s pollsvc.Service) Set {

	//
	createEndpoint := MakeEndpointCreate(s)
	createEndpoint = middleware.SetUserID(createEndpoint)

	//
	deleteEndpoint := MakeEndpointDelete(s)
	deleteEndpoint = middleware.SetUserID(deleteEndpoint)

	//
	getManyEndpoint := MakeEndpointGetMany(s)
	getManyEndpoint = middleware.SetUserID(getManyEndpoint)

	//
	getEndpoint := MakeEndpointGet(s)
	getEndpoint = middleware.SetUserID(getEndpoint)

	return Set{
		Create:  createEndpoint,
		Delete:  deleteEndpoint,
		GetMany: getManyEndpoint,
		Get:     getEndpoint,
	}
}

// MakeEndpointCreate create a create poll endpoint
func MakeEndpointCreate(s pollsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		p := models.Poll{
			Type:        models.PTypeSingle,
			Description: req.Description,
			Option:      req.Options,
			Tags:        req.Tag,
		}

		code, err := s.Create(ctx, p)
		if err != nil {
			return pollie.Response{
				Err: err,
			}, nil
		}

		return pollie.Response{
			Status: "success",
			Data: map[string]interface{}{
				"link":    fmt.Sprintf("https://pollie.io/v/%s", code),
				"message": "poll created successfully",
			},
		}, nil
	}
}

// MakeEndpointDelete creates a delete poll endpoint
func MakeEndpointDelete(s pollsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)
		err := s.Delete(ctx, req.ID, req.UserID)
		if err != nil {
			return pollie.Response{Err: err}, nil
		}

		return pollie.Response{
			Status: "success",
			Data: map[string]interface{}{
				"message": "poll deleted successfully",
			},
		}, nil
	}
}

// MakeEndpointGetMany create a get many polls endpoint
func MakeEndpointGetMany(s pollsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetManyReq)
		pp, err := s.GetMany(ctx, req.UserID)
		if err != nil {
			return pollie.Response{Err: err}, nil
		}

		return pollie.Response{
			Status: "success",
			Data: map[string]interface{}{
				"message": "successfully fetched polls",
				"polls":   pp,
				// TODO add pagination
			},
		}, nil
	}
}

// MakeEndpointGet create a get poll endpoint
func MakeEndpointGet(s pollsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		// get the poll without the user
		p, err := s.Get(ctx, req.ID)
		if err != nil {
			return pollie.Response{Err: err}, nil
		}

		// NOTE Auth does not seem to have a usecase at the moment
		if p.Auth {
			userID, _ := primitive.ObjectIDFromHex(middleware.GetUserID(ctx))
			if p.UserID != userID {
				return pollie.Response{Err: errors.New("user not authorized")}, nil
			}
		}

		return pollie.Response{
			Status: "success",
			Data: map[string]interface{}{
				"message": "successfully fetched poll",
				"poll":    p,
			},
		}, nil
	}
}
