package endpoints

import (
	"api-gateway/transports"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"net/url"
)

func MakeProxyIssueTokenEndpoint(proxyURL *url.URL) endpoint.Endpoint {
	return httptransport.NewClient(http.MethodGet,
		proxyURL,
		httptransport.EncodeJSONRequest,
		transports.DecodeIssueTokenResponse).Endpoint()
}
