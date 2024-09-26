package model

type SessionsReq struct {
	GrantType string `json:"grantType"`
	UserName  string `json:"userName"`
	Value     string `json:"value"`
}
