package model

import (
	"errors"
	"slices"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
)

/**
* system_id is the key for crm or other system connect with chat
 */
type ChatIntegrateSystem struct {
	*Base
	bun.BaseModel           `bun:"table:chat_integrate_system,alias:cis"`
	SystemName              string                    `json:"system_name" bun:"system_name,type:text,notnull"`
	Salt                    string                    `json:"salt" bun:"salt,type:text,notnull"`
	SystemId                string                    `json:"system_id" bun:"system_id,type:text,notnull"` // use for client connect header
	ChatAppIntegrateSystems []*ChatAppIntegrateSystem `json:"chat_app_integrate_systems" bun:"rel:has-many,join:id=chat_integrate_system_id"`
	TenantDefaultId         string                    `json:"tenant_default_id" bun:"tenant_default_id,type:uuid,notnull"`
	VendorId                string                    `json:"vendor_id" bun:"vendor_id,type:uuid,nullzero,default:null"`
	Vendor                  *ChatVendor               `json:"vendor" bun:"rel:belongs-to,join:vendor_id=id"`
	Status                  bool                      `json:"status" bun:"status,type:boolean,notnull"`
	InfoSystem              *InfoSystem               `json:"info_system" bun:"info_system,type:jsonb,notnull"`
}

type ChatIntegrateSystemView struct {
	*Base
	bun.BaseModel   `bun:"table:chat_integrate_system,alias:cis"`
	SystemName      string      `json:"system_name" bun:"system_name"`
	Salt            string      `json:"salt" bun:"salt-"`
	SystemId        string      `json:"system_id" bun:"system_id"` // use for client connect header
	TenantDefaultId string      `json:"tenant_default_id" bun:"tenant_default_id"`
	VendorId        string      `json:"vendor_id" bun:"vendor_id"`
	Vendor          *ChatVendor `json:"vendor" bun:"rel:belongs-to,join:vendor_id=id"`
	Status          bool        `json:"status" bun:"status"`
	InfoSystem      *InfoSystem `json:"info_system" bun:"info_system"`
}

type InfoSystem struct {
	AuthType            string `json:"auth_type"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	Token               string `json:"token"`
	WebsocketUrl        string `json:"websocket_url"`
	ApiUrl              string `json:"api_url"`
	ApiAuthUrl          string `json:"api_auth_url"`
	ApiGetUserDetailUrl string `json:"api_get_user_detail_url"`
	ApiGetUserUrl       string `json:"api_get_user_url"`
}

type ChatIntegrateSystemRequest struct {
	SystemName          string           `json:"system_name"`
	TenantDefaultId     string           `json:"tenant_default_id"`
	VendorId            string           `json:"vendor_id"`
	Status              bool             `json:"status"`
	AuthType            string           `json:"auth_type"`
	Username            string           `json:"username"`
	Password            string           `json:"password"`
	Token               string           `json:"token"`
	WebsocketUrl        string           `json:"websocket_url"`
	ApiUrl              string           `json:"api_url"`
	ApiAuthUrl          string           `json:"api_auth_url"`
	ApiGetUserDetailUrl string           `json:"api_get_user_detail_url"`
	ApiGetUserUrl       string           `json:"api_get_user_url"`
	ChatAppIds          []string         `json:"chat_app_ids"` // use for getting app_id available
	ChatApps            []ChatAppRequest `json:"chat_apps"`    // use for creating new app, this field and chat_app_ids are mutually exclusive
}

type AuthTypeJwtToken struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	AddTo string `json:"add_to"`
}

func (m *ChatIntegrateSystemRequest) Validate() error {
	if len(m.SystemName) < 1 {
		return errors.New("system name is required")
	}
	if len(m.VendorId) < 1 {
		return errors.New("vendor id is required")
	}
	if len(m.AuthType) < 1 {
		return errors.New("authorization is required")
	}
	if !slices.Contains(variables.AUTH, m.AuthType) {
		return errors.New("authorization " + m.AuthType + " is not supported")
	}
	switch m.AuthType {
	case "no_auth":
	case "basic_auth":
		if len(m.Username) < 1 {
			return errors.New("username is required")
		}
		if len(m.Password) < 1 {
			return errors.New("password is required")
		}
	case "bearer_token":
		if len(m.Token) < 1 {
			return errors.New("token is required")
		}
	case "jwt_token":
		if len(m.Token) < 1 {
			return errors.New("token is required")
		}
	case "pitel_crm":
	default:
	}

	if len(m.ApiUrl) < 1 {
		return errors.New("api url is required")
	}
	if len(m.ApiAuthUrl) < 1 {
		return errors.New("api auth url is required")
	}
	if len(m.ApiGetUserDetailUrl) < 1 {
		return errors.New("api get user detail url is required")
	}
	if len(m.ApiGetUserUrl) < 1 {
		return errors.New("api get user url is required")
	}
	return nil
}
