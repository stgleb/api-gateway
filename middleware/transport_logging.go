package middleware

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func LoggingMiddleware(logger log.Logger, endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", fmt.Sprintf("calling %s", endpointName))
			defer logger.Log("msg", fmt.Sprintf("called %s", endpointName))
			return next(ctx, request)
		}
	}
}
