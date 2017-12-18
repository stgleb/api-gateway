package middleware

import (
	"api-gateway"
	"context"
	"github.com/go-kit/kit/log"
)

type LoggingMiddleWare struct {
	log.Logger
	Next api_gateway.TokenService
}

func (mw LoggingMiddleWare) IssueToken(ctx context.Context, login, password string) (string, error) {
	mw.Logger.Log("method", "IssueToken", "login", login, "password", password)
	token, err := mw.Next.IssueToken(ctx, login, password)
	mw.Logger.Log("method", "IssueToken", token, err)

	return token, err
}

func (mw LoggingMiddleWare) VerifyToken(ctx context.Context, token string) error {
	mw.Logger.Log("method", "VerifyToken", "token", token)
	err := mw.Next.VerifyToken(ctx, token)
	mw.Logger.Log("method", "VerifyToken", "error", err.Error())

	return err
}

func (mw LoggingMiddleWare) RevokeToken(ctx context.Context, token string) error {
	mw.Logger.Log("method", "RevokeToken", "token", token)
	err := mw.Next.RevokeToken(ctx, token)
	mw.Logger.Log("method", "RevokeToken", "error", err.Error())

	return err
}
