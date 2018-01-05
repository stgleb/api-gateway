package endpoints

import (
	"api-gateway"
	. "api-gateway/data"
	"api-gateway/transports"
	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"net/url"
)

func MakeProxyIssueTokenEndpoint(proxyURL *url.URL) endpoint.Endpoint {
	return httptransport.NewClient(http.MethodPost,
		proxyURL,
		httptransport.EncodeJSONRequest,
		transports.DecodeIssueTokenResponse).Endpoint()
}

func MakeProxyVerifyTokenEndpoint(proxyURL *url.URL) endpoint.Endpoint {
	return httptransport.NewClient(http.MethodPost,
		proxyURL,
		httptransport.EncodeJSONRequest,
		transports.DecodeVerifyTokenResponse).Endpoint()
}

func MakeProxyRevokeTokenEndpoint(proxyURL *url.URL) endpoint.Endpoint {
	return httptransport.NewClient(http.MethodPost,
		proxyURL,
		httptransport.EncodeJSONRequest,
		transports.DecodeRevokeTokenResponse).Endpoint()
}

func MakeHealthCheckEndpoint(service api_gateway.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return HealthResponse{
			service.HealthCheck(),
		}, nil
	}
}
