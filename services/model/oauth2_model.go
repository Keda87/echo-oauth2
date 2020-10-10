package model

type RequestTokenPayload struct {
	GrantType    string `json:"grant_type" form:"grant_type"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	Scope        string `json:"scope,omitempty" form:"scope"`
}
