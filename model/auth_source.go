package model

import "github.com/uptrace/bun"

type AuthSource struct {
	bun.BaseModel `bun:"table:auth_source"`
	*Base
	TenantId string `json:"tenant_id"`
	Source   string `json:"source"`
	AuthUrl  string `json:"auth_url"`
	Info     *Info  `json:"info"`
	Status   bool   `json:"status"`
}

type Info struct {
	InfoType string `json:"info_type"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Authen struct {
	DomainUuid   string `json:"domain_uuid"`
	UserId       string `json:"user_id"`
	ClientId     string `json:"client_id"`
	ExpiredIn    int64  `json:"expired_in"`
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
	TokenType    string `json:"token_type"`
}

type AuthUserInfo struct {
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Username      string `json:"username" bun:"username"`
	Password      string `json:"password" bun:"password"`
	ApiKey        string `json:"api_key" bun:"api_key"`
	UserEnabled   string `json:"user_enabled" bun:"user_enabled"`
	UserStatus    string `json:"user_status" bun:"user_status"`
	Level         string `json:"level" bun:"level"`
	LastName      string `json:"last_name" bun:"last_name"`
	MiddleName    string `json:"middle_name" bun:"middle_name"`
	FirstName     string `json:"first_name" bun:"first_name"`
	UnitUuid      string `json:"unit_uuid" bun:"unit_uuid"`
	UnitName      string `json:"unit_name" bun:"unit_name"`
	RoleUuid      string `json:"role_uuid" bun:"role_uuid"`
	Extension     string `json:"extension" bun:"extension"`
	ExtensionUuid string `json:"extension_uuid" bun:"extension_uuid"`
}
