package pkg

import (
	"pollie/pkg/poll/pollendpoint"
	"pollie/pkg/poll/polltransport"
	"pollie/pkg/vote/voteendpoint"
	"pollie/pkg/vote/votetransport"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

// NewRouter create a router
func NewRouter(poll pollendpoint.Set, vote voteendpoint.Set) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger) //TODO remove
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Mount("/poll", polltransport.NewHandler(poll))
	r.Mount("/vote", votetransport.NewHandler(vote))

	return r
}
