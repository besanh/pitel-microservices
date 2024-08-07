package model

import (
	"time"

	"github.com/shaj13/go-guardian/v2/auth"
)

type (
	LoginBody struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
		Password string `json:"password" binding:"required,min=6,max=50"`
	}

	ForgotPassword struct {
		State     string
		Email     string
		Code      string
		Password  string
		Token     string
		UserAgent string
	}

	ForgotPasswordResponse struct {
		Code            string        `json:"code"`
		IsValid         bool          `json:"is_valid"`
		Email           string        `json:"email"`
		NextRequestTime time.Time     `json:"next_request_time"`
		RemainTime      time.Duration `json:"remain_time"`
	}
)

type BSS_DB struct {
	DBName   string `json:"database_name"`
	Port     int    `json:"database_port"`
	Host     string `json:"database_host"`
	Username string `json:"database_user"`
	Password string `json:"database_password"`
}

type LoginRequest struct {
	Username    string `json:"username" required:"true" binding:"required,min=6,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	UserAgent   string `json:"-"`
	Fingerprint string `json:"-"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" required:"true" minLength:"5"`
	Token        string `json:"token" required:"false"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	ExpiredIn    int    `json:"expired_in"`
	Token        string `json:"token"`
}

type LoginResponse struct {
	*AuthUserData
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredIn    int    `json:"expired_in"`
	Fullname     string `json:"fullname"`
	TenantName   string `json:"tenant_name"`
}

type AuthUser struct {
	*auth.DefaultUser
	*AuthUserData
	DB *BSS_DB
}

type AuthUserData struct {
	TenantId       string   `json:"tenant_id"`
	BusinessUnitId string   `json:"business_unit_id"`
	UserId         string   `json:"user_id"`
	Username       string   `json:"username"`
	Level          string   `json:"level"`
	Scopes         []string `json:"scopes"`
	RoleId         string   `json:"role_id"`
}

func (d *AuthUser) GetID() string {
	return d.UserId
}

func (d *AuthUser) GetLevel() string {
	return d.Level
}

func (d *AuthUser) SetLevel(level string) {
	d.Level = level
}

func (d *AuthUser) GetScopes() []string {
	return d.Scopes
}

func (d *AuthUser) SetScopes(scopes []string) {
	d.Scopes = scopes
}

func (d *AuthUser) IsLevel(levels ...string) bool {
	for _, l := range levels {
		if l == d.Level {
			return true
		}
	}
	return false
}
