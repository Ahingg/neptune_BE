package log_on

type LogOnResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // Seconds until access token expires
	RefreshToken string `json:"refresh_token"`
}
