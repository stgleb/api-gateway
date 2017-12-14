package api_gateway

import "context"

type TokenService interface {
	IssueToken(context.Context, string, string) (string, error)
	VerifyToken(context.Context, string) error
	RevokeToken(context.Context, string) error
}
