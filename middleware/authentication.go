package middleware

import (
	"context"
	"pollie"

	"github.com/go-chi/jwtauth"

	"github.com/go-kit/kit/endpoint"
)

var (
	userIDContext pollie.ContextKey = "user_id"
)

// SetUserID set the user id
func SetUserID(e endpoint.Endpoint) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t, cl, err := jwtauth.FromContext(ctx)

		var _ctx = ctx
		if err == nil && t.Valid {
			// check if uid exist in the claim
			uid, ok := cl["uid"].(string)

			if ok {
				_ctx = context.WithValue(ctx, userIDContext, uid)
			}
		}

		return e(_ctx, request)
	}

}

// GetUserID get the user id from a context
func GetUserID(ctx context.Context) string {
	v, ok := ctx.Value(userIDContext).(string)
	if ok {
		return v
	}

	return ""
}
