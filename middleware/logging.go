package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Logger log
func Logger(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (res interface{}, err error) {
			t := time.Now()

			defer func() {
				logger.Log(
					"level", "endpoint",
					"duration", time.Now().Sub(t),
					"time", time.Now(),
				)
			}()

			return next(ctx, request)
		}
	}

}
