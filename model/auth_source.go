package model

import "github.com/uptrace/bun"

type AuthSource struct {
	bun.BaseModel `bun:"table:auth_source"`
	*Base
	Source  string `json:"source"`
	AuthUrl string `json:"auth_url"`
	Info    *Info  `json:"info"`
	Status  bool   `json:"status"`
}

type Info struct {
	InfoType string `json:"info_type"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Authen struct {
	DomainUuid string `json:"domain_uuid"`
	UserUuid   string `json:"user_uuid"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	Extension  string `json:"extension"`
	Level      string `json:"level"`
}
