package model

import (
	"errors"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type ChatConnectionApp struct {
	*Base
	bun.BaseModel  `bun:"table:chat_connection_app,alias:cca"`
	TenantId       string `json:"tenant_id"`
	BusinessUnitId string `json:"business_unit_id"`
	ConnectionName string `json:"connection_name" bun:"connection_name,type:text,notnull"`
	ConnectionType string `json:"connection_type" bun:"connection_type,type:text,notnull"`
	AppId          string `json:"app_id" bun:"app_id,type:text,notnull"`
	QueueId        string `json:"queue_id" bun:"queue_id,type:text,notnull"`
	OaInfo         OaInfo `json:"oa_info" bun:"oa_info,type:text,notnull"`
	Status         string `json:"status" bun:"status,notnull"`
}

type OaInfo struct {
	Zalo []struct {
		UrlOa string `json:"url_oa"`
	} `json:"zalo"`
	Facebook []struct {
		OaId        string `json:"oa_id"`
		OaName      string `json:"oa_name"`
		AccessToken string `json:"access_token"`
	} `json:"facebook"`
}

type ChatConnectionAppRequest struct {
	ConnectionName string `json:"connection_name"`
	ConnectionType string `json:"connection_type"`
	QueueId        string `json:"queue_id"`
	OaInfo         OaInfo `json:"oa_info"`
	Status         string `json:"status"`
}

type AccessInfo struct {
	CallbackUrl   string `json:"callback_url"`
	ChallangeCode string `json:"challange_code"`
	State         string `json:"state"`
}

func (m *ChatConnectionAppRequest) Validate() error {
	if len(m.ConnectionName) < 1 {
		return errors.New("connection name is required")
	}
	if len(m.ConnectionType) < 1 {
		return errors.New("connection type is required")
	}
	if !slices.Contains[[]string](variables.CONNECTION_TYPE, m.ConnectionType) {
		return errors.New("connection type " + m.ConnectionType + " is not supported")
	}
	if len(m.Status) < 1 {
		return errors.New("status is required")
	}
	return nil
}
