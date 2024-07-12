package model

type AAAResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    *AuthUser `json:"data"`
}

type BssAuthRequest struct {
	Token   string `json:"token"`
	AuthUrl string `json:"auth_url"`
	Source  string `json:"source"`
	UserId  string `json:"user_id"`
}
