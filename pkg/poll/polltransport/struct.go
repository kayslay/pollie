package polltransport

import (
	"net/http"
	"pollie/pkg/poll/pollendpoint"
	"strings"

	"github.com/go-chi/jwtauth"
)

// implements a bind method to handle validation for CreateReq
type httpCreateReq pollendpoint.CreateReq

func (c *httpCreateReq) Bind(r *http.Request) error {
	c.Description = strings.TrimSpace(c.Description)

	_, cl, _ := jwtauth.FromContext(r.Context())

	if c.Auth {
		// token must be valid if poll is set to auth
		if err := cl.Valid(); err != nil {
			return err
		}
	}

	return (*pollendpoint.CreateReq)(c).Validate()
}
