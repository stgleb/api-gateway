package transports

import (
	. "api-gateway/data"
	"context"
	"encoding/json"
	"net/http"
)

func DecodeIssueTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var issueTokenRequest IssueTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&issueTokenRequest); err != nil {
		return nil, err
	}

	return issueTokenRequest, nil
}

func DecodeVerifyTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var verifyTokenRequest VerifyTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&verifyTokenRequest); err != nil {
		return nil, err
	}

	return verifyTokenRequest, nil
}

func DecodeRevokeTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var revokeTokenRequest RevokeTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&revokeTokenRequest); err != nil {
		return nil, err
	}

	return revokeTokenRequest, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeIssueTokenResponse(_ context.Context, r *http.Response) (response interface{}, err error) {
	var issueTokenResponse IssueTokenResponse

	if err := json.NewDecoder(r.Body).Decode(&issueTokenResponse); err != nil {
		return nil, err
	}

	return issueTokenResponse, nil
}

func DecodeVerifyTokenResponse(_ context.Context, r *http.Response) (response interface{}, err error) {
	var verifyTokenResponse VerifyTokenResponse

	if err := json.NewDecoder(r.Body).Decode(&verifyTokenResponse); err != nil {
		return nil, err
	}

	return verifyTokenResponse, nil
}

func DecodeRevokeTokenResponse(_ context.Context, r *http.Response) (response interface{}, err error) {
	var revokeTokenResponse RevokeTokenResponse

	if err := json.NewDecoder(r.Body).Decode(&revokeTokenResponse); err != nil {
		return nil, err
	}

	return revokeTokenResponse, nil
}
