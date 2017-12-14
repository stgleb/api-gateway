package middleware

import (
	"api-gateway"
	"api-gateway/data"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

type ProxyMiddleware struct {
	IssueTokenEndpoint  endpoint.Endpoint
	VerifyTokenEndpoint endpoint.Endpoint
	RevokeTokenEndpoint endpoint.Endpoint
}

type Middleware func(api_gateway.TokenService) api_gateway.TokenService

func (proxy ProxyMiddleware) IssueToken(ctx context.Context, login, password string) (string, error) {
	r, err := proxy.IssueTokenEndpoint(ctx, data.IssueTokenRequest{
		login,
		password,
	})

	if err != nil {
		return "", err
	}

	resp, ok := r.(data.IssueTokenResponse)

	if !ok {
		return "", errors.New(fmt.Sprintf("Error while converting response %v to IssueTokenResponse", r))
	}

	if len(resp.Error) > 0 {
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}

func (proxy ProxyMiddleware) VerifyToken(ctx context.Context, token string) error {
	r, err := proxy.VerifyTokenEndpoint(ctx, data.VerifyTokenRequest{
		token,
	})

	if err != nil {
		return err
	}

	resp, ok := r.(data.VerifyTokenResponse)

	if !ok {
		return errors.New(fmt.Sprintf("Error while converting response %v to VerifyTokenResponse", r))
	}

	if len(resp.Error) > 0 {
		return errors.New(resp.Error)
	}

	return nil
}

func (proxy ProxyMiddleware) RevokeToken(ctx context.Context, token string) error {
	r, err := proxy.RevokeTokenEndpoint(ctx, data.RevokeTokenRequest{
		token,
	})

	if err != nil {
		return err
	}

	resp, ok := r.(data.RevokeTokenResponse)

	if !ok {
		return errors.New(fmt.Sprintf("Error while converting response %v to RevokeTokenResponse", r))
	}

	if len(resp.Error) > 0 {
		return errors.New(resp.Error)
	}

	return nil
}
