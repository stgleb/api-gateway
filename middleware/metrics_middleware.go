package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"time"
)

func MetricsMiddleware(count metrics.Counter, latency metrics.Histogram, endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				count.With(endpointName).Add(1)
				latency.With(endpointName).Observe(time.Since(begin).Seconds())
			}(time.Now())

			return next(ctx, request)
		}
	}
}
