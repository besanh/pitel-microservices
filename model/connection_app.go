package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ConnectionApp struct {
	*Base
	bun.BaseModel  `bun:"table:connection_app,alias:ca"`
	ConnectionName string   `json:"connection_name" bun:"connection_name,type:text,notnull"`
	ConnectionType string   `json:"connection_type" bun:"connection_type,type:text,notnull"`
	TenantId       string   `json:"tenant_id"`
	BusinessUnitId string   `json:"business_unit_id"`
	AppId          string   `json:"app_id" bun:"app_id,type:text,notnull"`
	ChatApp        *ChatApp `json:"chat_app" bun:"chat_app,type:jsonb,notnull"`
	Status         bool     `json:"status" bun:"status,notnull"`
}

type ConnectionAppRequest struct {
	ConnectionName string `json:"connection_name"`
	ConnectionType string `json:"connection_type"`
	AppId          string `json:"app_id"`
	Status         string `json:"status"`
}

type AccessInfo struct {
	CallbackUrl   string `json:"callback_url"`
	ChallangeCode string `json:"challange_code"`
	State         string `json:"state"`
}

func (m *ConnectionAppRequest) Validate() error {
	if len(m.ConnectionName) < 1 {
		return errors.New("connection name is required")
	}
	if len(m.ConnectionType) < 1 {
		return errors.New("connection type is required")
	}
	if len(m.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(m.Status) > 0 {
		return errors.New("status is required")
	}
	return nil
}
