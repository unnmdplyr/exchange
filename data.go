package main

type Headers []map[string]string

type Body interface{}

type ServiceData struct {
	Method  string  `json:"method"`
	Url     string  `json:"url"`
	Headers Headers `json:"headers"`
	Body    Body    `json:"body"`
}

type EnvironmentData map[string]string

type TokenData struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	Token        string `json:"token"`
}
