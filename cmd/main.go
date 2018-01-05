package main

import (
	. "api-gateway"
	. "api-gateway/endpoints"
	. "api-gateway/middleware"
	. "api-gateway/services"
	. "api-gateway/transports"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	. "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"os"

	"bufio"
	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configFileName string

	issueTokenCounter  metrics.Counter
	verifyTokenCounter metrics.Counter
	revokeTokenCounter metrics.Counter

	issueTokenHistogram  metrics.Histogram
	verifyTokenHistogram metrics.Histogram
	revokeTokenHistogram metrics.Histogram

	issueTokenLabel  = "issueToken"
	verifyTokenLabel = "verifyToken"
	revokeTokenLabel = "revokeToken"
)

func init() {
	flag.StringVar(&configFileName, "config", "config.toml", "config file name")
	flag.Parse()
}

func main() {
	config, err := NewConfig(configFileName)

	if err != nil {
		panic(err)
	}

	logfile, err := os.OpenFile(config.Main.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	bufferedFileWriter := bufio.NewWriter(logfile)

	defer logfile.Close()

	fileWriter := log.NewSyncWriter(bufferedFileWriter)
	logWriter := io.MultiWriter(fileWriter, os.Stderr)

	logger := log.NewLogfmtLogger(logWriter)

	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var tokenService TokenService

	if err != nil {
		logger.Log(err)
	}

	issueTokenProxyURL := &url.URL{
		Scheme: config.TokenService.Protocol,
		Host:   config.TokenService.ListenStr,
		Path:   config.TokenService.IssueTokenPath,
	}

	verifyTokenProxyURL := &url.URL{
		Scheme: config.TokenService.Protocol,
		Host:   config.TokenService.ListenStr,
		Path:   config.TokenService.VerifyTokenPath,
	}

	revokeTokenProxyURL := &url.URL{
		Scheme: config.TokenService.Protocol,
		Host:   config.TokenService.ListenStr,
		Path:   config.TokenService.RevokeTokenPath,
	}

	issueTokenEndpoint := MakeProxyIssueTokenEndpoint(issueTokenProxyURL)
	verifyTokenEndpoint := MakeProxyVerifyTokenEndpoint(verifyTokenProxyURL)
	revokeTokenEndpoint := MakeProxyRevokeTokenEndpoint(revokeTokenProxyURL)

	issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint =
		wrapLogging(issueTokenEndpoint, logger, verifyTokenEndpoint, revokeTokenEndpoint)
	issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint =
		wrapLimit(issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint)
	issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint =
		wrapPrometheus(config, issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint)

	tokenService = TokenProxyService{
		issueTokenEndpoint,
		verifyTokenEndpoint,
		revokeTokenEndpoint,
	}

	tokenService = NewLoggingMiddleWare(tokenService, logger)

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

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/token", issueTokenHandler)
	http.Handle("/token/verify", verifyTokenHandler)
	http.Handle("/token/revoke", revokerTokenHandler)

	logger.Log("main", fmt.Sprintf("Start listen port %s", config.TokenService.ListenStr))
	logger.Log("error", http.ListenAndServe(config.TokenService.ListenStr, nil))
}

func wrapLogging(issueTokenEndpoint endpoint.Endpoint, logger log.Logger, verifyTokenEndpoint endpoint.Endpoint,
	revokeTokenEndpoint endpoint.Endpoint) (endpoint.Endpoint, endpoint.Endpoint, endpoint.Endpoint) {
	issueTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "IssueToken"),
		"issueTokenEndpoint")(issueTokenEndpoint)
	verifyTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "VerifyToken"),
		"verifyTokenEndpoint")(verifyTokenEndpoint)
	revokeTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "RevokeToken"),
		"revokeTokenEndpoint")(revokeTokenEndpoint)

	return issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint
}

func wrapLimit(issueTokenEndpoint endpoint.Endpoint, verifyTokenEndpoint endpoint.Endpoint,
	revokeTokenEndpoint endpoint.Endpoint) (endpoint.Endpoint, endpoint.Endpoint, endpoint.Endpoint) {
	rateLimitMiddleware10 := NewErroringLimiter(rate.NewLimiter(10, 1))
	rateLimitMiddleware5 := NewErroringLimiter(rate.NewLimiter(5, 1))
	rateLimitMiddleware1 := NewErroringLimiter(rate.NewLimiter(1, 1))

	issueTokenEndpoint = rateLimitMiddleware10(issueTokenEndpoint)
	verifyTokenEndpoint = rateLimitMiddleware5(verifyTokenEndpoint)
	revokeTokenEndpoint = rateLimitMiddleware1(revokeTokenEndpoint)

	return issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint
}

func wrapPrometheus(config *TomlConfig, issueTokenEndpoint endpoint.Endpoint, verifyTokenEndpoint endpoint.Endpoint,
	revokeTokenEndpoint endpoint.Endpoint) (endpoint.Endpoint, endpoint.Endpoint, endpoint.Endpoint) {
	// Issue token
	issueTokenCounter := kitprometheus.NewCounterFrom(
		stdprometheus.CounterOpts{
			Name: "issue_token_counter",
			Subsystem: config.Main.ServiceName,
			Help: "Issue token counter"},
		[]string{issueTokenLabel})
	issueTokenHistogram := kitprometheus.NewHistogramFrom(
		stdprometheus.HistogramOpts{
			Name: "issue_token_histogram",
			Subsystem: config.Main.ServiceName,
			Help: "Issue token histogram"},
		[]string{issueTokenLabel})
	issueTokenEndpoint = MetricsMiddleware(issueTokenCounter, issueTokenHistogram, issueTokenLabel)(issueTokenEndpoint)

	// Verify token
	verifyTokenCounter = kitprometheus.NewCounterFrom(
		stdprometheus.CounterOpts{
			Name: "verify_token_counter",
			Subsystem: config.Main.ServiceName,
			Help: "Verify token counter",
		}, []string{verifyTokenLabel})
	verifyTokenHistogram = kitprometheus.NewHistogramFrom(
		stdprometheus.HistogramOpts{
			Name: "verify_token_histogram",
			Subsystem: config.Main.ServiceName,
			Help: "Verify token histogram",
		}, []string{verifyTokenLabel})
	verifyTokenEndpoint = MetricsMiddleware(verifyTokenCounter, verifyTokenHistogram, verifyTokenLabel)(verifyTokenEndpoint)

	// Revoke token
	revokeTokenCounter = kitprometheus.NewCounterFrom(
		stdprometheus.CounterOpts{
			Name: "revoke_token_counter",
			Subsystem: config.Main.ServiceName,
			Help: "Revoke token counter",
		},
		[]string{revokeTokenLabel})
	revokeTokenHistogram = kitprometheus.NewHistogramFrom(
		stdprometheus.HistogramOpts{
			Name: "revoke_token_histogram",
			Subsystem: config.Main.ServiceName,
			Help: "Revoke token histogram",
		},
		[]string{revokeTokenLabel})

	revokeTokenEndpoint = MetricsMiddleware(revokeTokenCounter, revokeTokenHistogram, revokeTokenLabel)(revokeTokenEndpoint)

	return issueTokenEndpoint, verifyTokenEndpoint, revokeTokenEndpoint
}
