package polltransport

import (
	"context"
	"fmt"
	"net/http"
	"pollie"
	"pollie/middleware"
	"pollie/pkg/poll/pollendpoint"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/spf13/viper"
)

// NewHandler create a new Handler
func NewHandler(s pollendpoint.Set) http.Handler {

	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(pollie.ErrorEncoder),
		// httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	createHandler := httptransport.NewServer(
		s.Create,
		decodeCreateRequest,
		pollie.EncodeResponse,
		options...,
	)

	deleteHandler := httptransport.NewServer(
		s.Delete,
		decodeDeleteRequest,
		pollie.EncodeResponse,
		options...,
	)
	getManyHandler := httptransport.NewServer(
		s.GetMany,
		decodeGetManyRequest,
		pollie.EncodeResponse,
		options...,
	)
	getHandler := httptransport.NewServer(
		s.Get,
		decodeGetRequest,
		pollie.EncodeResponse,
		options...,
	)

	r.Route("/", func(r chi.Router) {
		var tokenAuth = jwtauth.New("HS256", []byte(viper.GetString("JWT_SECRET")), nil)

		r.Use(jwtauth.Verifier(tokenAuth))

		r.Get("/", pollie.Handler2HandlerFunc(getManyHandler))
		r.Get("/{id}", pollie.Handler2HandlerFunc(getHandler))
		r.Post("/create", pollie.Handler2HandlerFunc(createHandler))
		r.Delete("/delete/{id}", pollie.Handler2HandlerFunc(deleteHandler))
	})

	return r
}

// decodeCreateRequest decode create request body
func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := httpCreateReq{}
	if err := render.Bind(r, &data); err != nil {
		return data, err
	}

	return pollendpoint.CreateReq(data), nil
}

// decodeDeleteRequest decode delete request
func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID := middleware.GetUserID(r.Context())

	data := pollendpoint.DeleteReq{
		ID:     chi.URLParam(r, "id"),
		UserID: userID,
	}
	return data, nil
}

// decodeGetManyRequest decode delete request
func decodeGetManyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID := middleware.GetUserID(r.Context())
	query := r.URL.Query()

	q := query.Get("q")

	data := pollendpoint.GetManyReq{
		UserID: userID,
		Q:      q,
	}

	return data, nil
}

// decodeGetRequest decode delete request
func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID := middleware.GetUserID(r.Context())
	fmt.Println(userID)

	data := pollendpoint.GetReq{
		ID:     chi.URLParam(r, "id"),
		UserID: userID,
	}

	return data, nil
}
