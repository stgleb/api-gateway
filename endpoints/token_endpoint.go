package endpoints

import (
	. "api-gateway"
	. "api-gateway/data"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

func MakeIssueTokenEndpoint(srv TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(IssueTokenRequest)

		if !ok {
			return nil, errors.New("wrong request format expected IssueTokenRequest")
		}

		r, err := srv.IssueToken(ctx, req.Login, req.Password)

		var errStr string

		if err != nil {
			errStr = err.Error()
		}

		return IssueTokenResponse{
			TokenResponse{
				r,
				errStr,
			},
		}, nil
	}
}

func MakeVerifyTokenEndpoint(srv TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(VerifyTokenRequest)

		if !ok {
			return nil, errors.New("wrong request format expected VerifyTokenRequest")
		}

		err := srv.VerifyToken(ctx, req.Token)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		return IssueTokenResponse{
			TokenResponse{
				"",
				errStr,
			},
		}, nil
	}
}

func MakeRevokeTokenEndpoint(srv TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RevokeTokenRequest)

		if !ok {
			return nil, errors.New("wrong request format expected RevokeTokenRequest")
		}

		err := srv.RevokeToken(ctx, req.Token)
		var errStr string

		if err != nil {
			errStr = err.Error()
		}

		return RevokeTokenResponse{
			TokenResponse{
				"",
				errStr,
			},
		}, nil
	}
}
