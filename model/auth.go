package model

type ChatResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    *AuthUser `json:"data"`
}

type BssAuthRequest struct {
	Token         string `json:"token"`
	ApiUrl        string `json:"api_url"`
	AuthUrl       string `json:"auth_url"`
	UserDetailUrl string `json:"user_detail_url"`
	Source        string `json:"source"`
	UserId        string `json:"user_id"`
	ServerName    string `json:"server_name"`
}
