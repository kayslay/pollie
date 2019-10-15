package votetransport

import (
	"context"
	"net/http"
	"pollie"
	"pollie/pkg/vote/voteendpoint"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/spf13/viper"
)

type httpVote voteendpoint.Vote

func (hv *httpVote) Bind(r *http.Request) error {

	return (*voteendpoint.Vote)(hv).Validate(chi.URLParam(r, "id"), strings.Split(r.RemoteAddr, ":")[0])
}

func NewHandler(s voteendpoint.Set) http.Handler {

	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(pollie.ErrorEncoder),
		// httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	voteHandler := httptransport.NewServer(
		s.Vote,
		decodeVoteRequest,
		pollie.EncodeResponse,
		options...,
	)

	r.Route("/", func(r chi.Router) {
		var tokenAuth = jwtauth.New("HS256", []byte(viper.GetString("JWT_SECRET")), nil)

		r.Use(jwtauth.Verifier(tokenAuth))

		r.Post("/{id}", pollie.Handler2HandlerFunc(voteHandler))
	})

	return r
}

func decodeVoteRequest(_ context.Context, r *http.Request) (interface{}, error) {

	data := httpVote{}

	if err := render.Bind(r, &data); err != nil {
		return data, err
	}

	data.Device = r.Header.Get("POLLIE-DEVICE")
	data.UserAgent = r.Header.Get("User-Agent")

	return voteendpoint.Vote(data), nil
}
