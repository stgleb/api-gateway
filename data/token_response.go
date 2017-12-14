package data

type TokenResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type IssueTokenResponse struct {
	TokenResponse
}

type VerifyTokenResponse struct {
	TokenResponse
}

type RevokeTokenResponse struct {
	TokenResponse
}
