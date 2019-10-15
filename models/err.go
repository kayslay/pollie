package models

import "errors"

var (
	// ErrInvalidID invalid objectID passed
	ErrInvalidID = errors.New("invalid id passed")
	// ErrInvalidPollType invalid poll type passed
	ErrInvalidPollType = errors.New("invalid poll type")
)
