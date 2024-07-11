package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatIntegrateSystem struct {
	*Base
	bun.BaseModel `bun:"table:chat_integrate_system,alias:cis"`
	SystemName    string      `json:"system_name" bun:"system_name,type:text,notnull"`
	VendorId      string      `json:"vendor_id" bun:"vendor_id,type:uuid,nullzero,default:null"`
	Vendor        *ChatVendor `json:"vendor" bun:"rel:belongs-to,join:vendor_id=id"`
	Status        bool        `json:"status" bun:"status,type:boolean,notnull"`
	InfoSystem    *InfoSystem `json:"info_system" bun:"info_system,type:jsonb,notnull"`
}

type InfoSystem struct {
	Authorization string `json:"authorization"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ApiKey        string `json:"api_key"`
	Token         string `json:"token"`
	WebsocketUrl  string `json:"websocket_url"`
	ApiUrl        string `json:"api_url"`
	ApiGetUserUrl string `json:"api_get_user_url"`
}

type ChatIntegrateSystemRequest struct {
	SystemName string      `json:"system_name"`
	VendorId   string      `json:"vendor_id"`
	Status     bool        `json:"status"`
	InfoSystem *InfoSystem `json:"info_system"`
}

func (m *ChatIntegrateSystemRequest) Validate() error {
	if len(m.SystemName) < 1 {
		return errors.New("system name is required")
	}
	if len(m.VendorId) < 1 {
		return errors.New("vendor id is required")
	}
	if m.InfoSystem == nil {
		return errors.New("info system is required")
	}
	if m.InfoSystem != nil {
		if len(m.InfoSystem.Authorization) < 1 {
			return errors.New("authorization is required")
		}
		if len(m.InfoSystem.Username) < 1 {
			return errors.New("username is required")
		}
		if len(m.InfoSystem.Password) < 1 {
			return errors.New("password is required")
		}
		if len(m.InfoSystem.ApiKey) < 1 {
			return errors.New("api key is required")
		}
		if len(m.InfoSystem.Token) < 1 {
			return errors.New("token is required")
		}
		if len(m.InfoSystem.WebsocketUrl) < 1 {
			return errors.New("websocket url is required")
		}
		if len(m.InfoSystem.ApiUrl) < 1 {
			return errors.New("api url is required")
		}
		if len(m.InfoSystem.ApiGetUserUrl) < 1 {
			return errors.New("api get user url is required")
		}
	}
	return nil
}
