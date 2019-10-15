package pollie

import (
	"context"
	"encoding/json"
	"net/http"
)

type ContextKey string

type Failer interface {
	Failed() error
}

// Response struct for response
type Response struct {
	Status string                 `json:"status" bson:"status"`
	Data   map[string]interface{} `json:"data,omitempty," bson:"data"`
	Err    error                  `json:"-" bson:"err"`
}

// Failed returns an error if one exist
func (r Response) Failed() error {
	return r.Err
}

type errorWrapper struct {
	Status string
	Error  string
}

// ErrorEncoder
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error(), Status: "error"})
}

func err2code(err error) int {

	switch err.(type) {
	case LoggableErr:
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}

func Handler2HandlerFunc(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

// encodeResponse encode the response
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	w.Header().Add("Content-Type", "application/json")

	if f, ok := response.(Failer); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}
