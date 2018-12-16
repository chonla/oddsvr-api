package run

type InvertedToken struct {
	ID uint32 `json:"_id"`
	*Token
}

type Token struct {
	TokenType    string `json:"token_type"`
	Expiry       int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      `json:"athlete"`
}
