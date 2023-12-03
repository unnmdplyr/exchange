package main

type Headers []map[string]string
type Body map[string]string

type ServiceData struct {
	Method  string  `json:"method"`
	Url     string  `json:"url"`
	Headers Headers `json:"headers"`
	Body    Body    `json:"body"`
	// TODO: create a data type which can deal with compound objects
}

type EnvironmentData map[string]string

type TokenData struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}
