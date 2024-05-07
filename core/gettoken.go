package core

import (
	"encoding/json"
	"net/url"
)

// AuthInfo type contains the information needed for authentication
type AuthInfo struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Token type contains the token information of the current session.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}

// GetToken function returns a token based on the AuthInfo
func (ai AuthInfo) GetToken(baseURL string) Token {
	headers := map[string][]string{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Accept":       []string{"application/json"}}
	endpoint := "/SASLogon/oauth/token"
	method := "POST"
	data := url.Values{}
	// data.Set("username", ai.Username)
	// data.Set("password", ai.Password)
	data.Set("code", ai.Code)
	data.Set("grant_type", ai.GrantType)
	data.Set("client_id", ai.ClientID)
	data.Set("client_secret", ai.ClientSecret)
	// log.Println("Before CallRest")
	resp := CallRest(baseURL, endpoint, headers, method, data, nil)
	var token Token
	json.Unmarshal(resp, &token)
	return token
}
