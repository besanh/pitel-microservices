package model

import (
	"github.com/uptrace/bun"
)

type (
	IBKUser struct {
		*Base
		bun.BaseModel  `bun:"table:ibk_users,alias:u"`
		BusinessUnitId string   `bun:"business_unit_id,type:uuid,notnull" json:"business_unit_id"`
		Username       string   `bun:"username,type:text,notnull" json:"username"`
		Fullname       string   `bun:"fullname,type:text" json:"fullname"`
		Email          string   `bun:"email,type:text" json:"email"`
		PhoneNumber    string   `bun:"phone_number,type:text" json:"phone_number"`
		Password       string   `bun:"password,type:text,notnull" json:"password"`
		Level          string   `bun:"level,type:text" json:"level"`
		Scopes         []string `bun:"scopes,array" json:"scopes"`
		IsActivated    bool     `bun:"is_activated,type:bool" json:"is_activated"`
		IsLocked       bool     `bun:"is_locked,type:bool" json:"is_locked"`
		IsSentEmail    bool     `bun:"is_sent_email,type:bool" json:"is_sent_email"`
		RoleId         string   `bun:"role_id,type:uuid,nullzero" json:"role_id"`
	}
	IBKUserInfo struct {
		*IBKUser
		bun.BaseModel    `bun:"table:ibk_users,alias:u"`
		BusinessUnitId   string `bun:"business_unit_id" json:"business_unit_id"`
		BusinessUnitName string `bun:"business_unit_name" json:"business_unit_name"`
		TenantId         string `bun:"tenant_id" json:"tenant_id"`
		RoleId           string `bun:"role_id" json:"role_id"`
		RoleName         string `bun:"role_name" json:"role_name"`
	}

	IBKUserShortInfo struct {
		bun.BaseModel `bun:"table:ibk_users,alias:u"`
		Id            string `bun:"id,pk" json:"id"`
		UserCode      string `bun:"user_code" json:"user_code"`
		Username      string `bun:"username" json:"username"`
		Fullname      string `bun:"fullname" json:"fullname"`
	}
)

type (
	IBKUserQueryParam struct {
		Keyword           string `query:"keyword"`
		Sort              string `query:"sort"`
		Order             string `query:"order"`
		BusinessUnitId_Eq string `query:"business_unit_id_eq"`
		TenantId_Eq       string `query:"tenant_id_eq"`
		RoleId_Eq         string `query:"role_id_eq"`
		ServiceId_Eq      string `query:"service_id_eq"`
		Username_Eq       string `query:"username_eq"`
		Fullname_Eq       string `query:"fullname_eq"`
		Email_Eq          string `query:"email_eq"`
		PhoneNumber_Eq    string `query:"phone_number_eq"`
		PhoneNumber_Like  string `query:"phone_number_like"`
		Fullname_Like     string `query:"fullname_like"`
		Email_Like        string `query:"email_like"`
		IsActivated_Eq    bool   `query:"is_activated"`
		IsLocked_Eq       bool   `query:"is_locked"`
	}
	IBKUserBody struct {
		_              struct{} `json:"-" additionalProperties:"true"`
		Username       string   `json:"username" required:"true" pattern:"^[a-zA-Z0-9_-]{2,50}$" patternDescription:"alphanumeric characters only" doc:"Username"`
		Password       string   `json:"password" required:"true" minLength:"8" maxLength:"50" pattern:"^[a-zA-Z0-9!@#$%^&*]{0,50}$" doc:"Password"`
		Fullname       string   `json:"fullname" required:"true" pattern:"^[a-zA-Z0-9 _-]{0,50}$" doc:"Fullname"`
		BusinessUnitId string   `json:"business_unit_id" required:"false" nullable:"true" format:"uuid" doc:"Business Unit Id"`
		Email          string   `json:"email" required:"false" format:"email" maxLength:"50" doc:"Email"`
		PhoneNumber    string   `json:"phone_number" required:"false" doc:"Phone number"`
		IsActivated    bool     `json:"is_activated" required:"false" default:"false" doc:"Is activated"`
		IsLocked       bool     `json:"is_locked" required:"false" default:"false" doc:"Is locked"`
		IsSentEmail    bool     `json:"is_sent_email" required:"false" default:"false" doc:"Is sent email"`
		Level          string   `json:"level" required:"true" enum:"user,admin,superadmin" doc:"Level"`
		Scopes         []string `json:"scopes" required:"false" doc:"Scopes"`
		RoleId         string   `json:"role_id" required:"true" nullable:"true" format:"uuid" doc:"Role Id"`
	}
)
