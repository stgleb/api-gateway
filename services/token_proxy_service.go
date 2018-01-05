package services

import (
	"api-gateway/data"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

type TokenProxyService struct {
	IssueTokenEndpoint  endpoint.Endpoint
	VerifyTokenEndpoint endpoint.Endpoint
	RevokeTokenEndpoint endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (proxy TokenProxyService) IssueToken(ctx context.Context, login, password string) (string, error) {
	r, err := proxy.IssueTokenEndpoint(ctx, data.LoginRequest{
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

func (proxy TokenProxyService) VerifyToken(ctx context.Context, token string) error {
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

func (proxy TokenProxyService) RevokeToken(ctx context.Context, token string) error {
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

func (proxy TokenProxyService) HealthCheck() bool {
	return true
}
