package data

type IssueTokenRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenRequest struct {
	Token string `json:"token"`
}

type VerifyTokenRequest struct {
	TokenRequest
}

type RevokeTokenRequest struct {
	TokenRequest
}
