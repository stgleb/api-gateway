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
	. "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"os"
	"io"
	"github.com/go-kit/kit/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	configFileName string

	issueTokenCounter metrics.Counter
	verifyTokenCounter metrics.Counter
	revokeTokenCounter metrics.Counter


	issueTokenHistogram metrics.Histogram
	verifyTokenHistogram metrics.Histogram
	revokeTokenHistogram metrics.Histogram

	issueTokenLabel = "issueToken"
	verifyTokenLabel = "verifyToken"
	revokeTokenLabel= "revokeToken"
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

	// TODO(stgleb): Make file writer buffered
	logfile, err := os.OpenFile(config.Main.LogFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	defer logfile.Close()

	fileWriter := log.NewSyncWriter(logfile)
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

	issueTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "IssueToken"),
		"issueTokenEndpoint")(issueTokenEndpoint)
	verifyTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "VerifyToken"),
		"verifyTokenEndpoint")(verifyTokenEndpoint)
	revokeTokenEndpoint = LoggingMiddleware(log.With(logger, "method", "RevokeToken"),
		"revokeTokenEndpoint")(revokeTokenEndpoint)

	// Cover issue token with limiter
	rateLimitMiddleware10 := NewErroringLimiter(rate.NewLimiter(10, 1))
	rateLimitMiddleware5 := NewErroringLimiter(rate.NewLimiter(5, 1))
	rateLimitMiddleware1 := NewErroringLimiter(rate.NewLimiter(1, 1))

	issueTokenEndpoint = rateLimitMiddleware10(issueTokenEndpoint)
	verifyTokenEndpoint = rateLimitMiddleware5(verifyTokenEndpoint)
	revokeTokenEndpoint = rateLimitMiddleware1(revokeTokenEndpoint)

	// Issue token
	issueTokenCounter := kitprometheus.NewCounter(stdprometheus.NewCounterVec(stdprometheus.CounterOpts{
		Name: issueTokenLabel,
	}, []string{issueTokenLabel}))

	issueTokenHistogram := kitprometheus.NewHistogram(stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
		Name: issueTokenLabel,
	}, []string{issueTokenLabel}))

	issueTokenEndpoint = MetricsMiddleware(issueTokenCounter, issueTokenHistogram, issueTokenLabel)(issueTokenEndpoint)

	// Verify token
	verifyTokenCounter = kitprometheus.NewCounter(stdprometheus.NewCounterVec(stdprometheus.CounterOpts{
		Name: verifyTokenLabel,
	}, []string{verifyTokenLabel}))

	verifyTokenHistogram = kitprometheus.NewHistogram(stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
		Name: verifyTokenLabel,
	}, []string{verifyTokenLabel}))

	verifyTokenEndpoint = MetricsMiddleware(verifyTokenCounter, verifyTokenHistogram, verifyTokenLabel)(verifyTokenEndpoint)

	// Revoke token
	revokeTokenCounter = kitprometheus.NewCounter(stdprometheus.NewCounterVec(stdprometheus.CounterOpts{
		Name: revokeTokenLabel,
	}, []string{revokeTokenLabel}))

	revokeTokenHistogram = kitprometheus.NewHistogram(stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
		Name: revokeTokenLabel,
	}, []string{revokeTokenLabel}))

	revokeTokenEndpoint = MetricsMiddleware(revokeTokenCounter, revokeTokenHistogram, revokeTokenLabel)(revokeTokenEndpoint)

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

	logger.Log("main", fmt.Sprintf("Start listen port %s", config.Main.ListenStr))
	logger.Log("error", http.ListenAndServe(config.Main.ListenStr, nil))
}
