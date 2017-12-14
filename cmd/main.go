package main

import (
	"golang.org/x/time/rate"
	. "api-gateway"
	. "api-gateway/endpoints"
	. "api-gateway/services"
	. "api-gateway/transports"
	. "api-gateway/middleware"
	"flag"
	"github.com/go-kit/kit/log"
	. "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

var (
	configFileName string
)

func init() {
	flag.StringVar(&configFileName, "config", "config.toml", "config file name")
	flag.Parse()
}

func main() {
	config, err := NewConfig(configFileName)
	logger := log.NewLogfmtLogger(os.Stderr)

	if err != nil {
		logger.Log(err)
	}

	var tokenService TokenService = TokenServiceImpl{}

	issueTokenEndpoint := MakeIssueTokenEndpoint(tokenService)
	verifyTokenEndpoint := MakeVerifyTokenEndpoint(tokenService)
	revokeTokenEndpoint := MakeRevokeTokenEndpoint(tokenService)

	issueTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "IssueToken"))(issueTokenEndpoint)
	verifyTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "VerifyToken"))(verifyTokenEndpoint)
	revokeTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "RevokeToken"))(revokeTokenEndpoint)

	// Cover issue token with limiter
	rateLimitMiddleware10 := NewErroringLimiter(rate.NewLimiter(10, 1))
	rateLimitMiddleware5 := NewErroringLimiter(rate.NewLimiter(5, 1))
	rateLimitMiddleware1 := NewErroringLimiter(rate.NewLimiter(1, 1))

	issueTokenEndpoint = rateLimitMiddleware10(issueTokenEndpoint)
	verifyTokenEndpoint = rateLimitMiddleware5(verifyTokenEndpoint)
	revokeTokenEndpoint = rateLimitMiddleware1(revokeTokenEndpoint)

	tokenService = LoggingMiddleWare{
		logger,
		tokenService,
	}

	issueTokenHandler := httptransport.NewServer(
		issueTokenEndpoint,
		DecodeIssueTokenRequest,
		EncodeResponse,
	)

	verifyTokenHandler := httptransport.NewServer(
		verifyTokenEndpoint,
		DecodeVerifyTokenRequest,
		EncodeResponse,
	)

	revokerTokenHandler := httptransport.NewServer(
		revokeTokenEndpoint,
		DecodeRevokeTokenRequest,
		EncodeResponse,
	)

	http.Handle("/token", issueTokenHandler)
	http.Handle("/token/verify", verifyTokenHandler)
	http.Handle("/token/revoke", revokerTokenHandler)

	logger.Log("Start listen port %s\n", config.Main.ListenStr)
	logger.Log(http.ListenAndServe(config.Main.ListenStr, nil))
}
