package model

type UserInfo struct {
	GrantType string `json:"grantType"`
	UserName  string `json:"userName"`
	Value     string `json:"value"`
}

type TokenInfo struct {
	AccessSession string `json:"accessSession"`
}
