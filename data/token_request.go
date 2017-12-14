package data

type IssueTokenRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

type RevokeTokenRequest struct {
	Token string `json:"token"`
}
