package transports

import (
	. "api-gateway/data"
	"context"
	"net/http"
)

// decode health check
func DecodeHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}
