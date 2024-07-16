package model

import "github.com/shaj13/go-guardian/v2/auth"

type AuthUser struct {
	TenantId          string          `json:"tenant_id"`
	UserId            string          `json:"user_id"`
	Username          string          `json:"username"`
	Level             string          `json:"level"`
	Source            string          `json:"source"`
	Token             string          `json:"token"`
	Fullname          string          `json:"fullname"`
	SecretKey         string          `json:"secret_key"`
	RoleId            string          `json:"role_id"`
	IntegrateSystemId string          `json:"integrate_system_id"`
	Extensions        auth.Extensions `json:"extensions"`
	Groups            []string        `json:"groups"`
	SystemId          string          `json:"system_id"`
	ApiUrl            string          `json:"api_url"`
	Server            string          `json:"server"`
}

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	UserAgent   string `json:"user_agent"`
	Fingerprint string `json:"-"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredIn    int    `json:"expired_in"`
	Fullname     string `json:"fullname"`
	Username     string `json:"username"`
	UserId       string `json:"user_id"`
	TenantLogo   string `json:"tenant_logo"`
	TenantId     string `json:"tenant_id"`
	TenantName   string `json:"tenant_name"`
	RoleId       string `json:"role_id"`
	Level        string `json:"level"`
	SystemId     string `json:"system_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	Token        string `json:"token"`
}

type TokenData struct {
	TenantId  string `json:"tenant_id"`
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	Level     string `json:"level"`
	Source    string `json:"source"`
	Token     string `json:"token"`
	Fullname  string `json:"fullname"`
	SecretKey string `json:"secret_key"`
	RoleId    string `json:"role_id"`
	SystemId  string `json:"system_id"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	ExpiredIn    int    `json:"expired_in"`
	Token        string `json:"token"`
}

func (m *AuthUser) GetUserName() string {
	return m.Username
}

func (m *AuthUser) SetUserName(username string) {
	m.Username = username
}

func (m *AuthUser) GetID() string {
	return m.UserId
}

func (m *AuthUser) SetID(id string) {
	m.UserId = id
}

func (m *AuthUser) GetLevel() string {
	return m.Level
}

func (m *AuthUser) SetLevel(level string) {
	m.Level = level
}

func (m *AuthUser) GetExtensions() auth.Extensions {
	if m.Extensions == nil {
		m.Extensions = auth.Extensions{}
	}
	return m.Extensions
}

func (m *AuthUser) SetExtensions(exts auth.Extensions) {
	m.Extensions = exts
}

func (m *AuthUser) GetGroups() []string {
	return m.Groups
}

func (m *AuthUser) SetGroups(groups []string) {
	m.Groups = groups
}

func (m *AuthUser) GetRoleID() string {
	return m.RoleId
}

func (m *AuthUser) SetRoleID(roleId string) {
	m.RoleId = roleId
}
