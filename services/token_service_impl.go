package services

import "context"

type TokenServiceImpl struct{}

func (tokenService TokenServiceImpl) IssueToken(ctx context.Context, login, password string) (string, error) {
	return "", nil
}

func (tokenService TokenServiceImpl) VerifyToken(ctx context.Context, token string) error {
	return nil
}

func (tokenService TokenServiceImpl) RevokeToken(ctx context.Context, token string) error {
	return nil
}
