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
	ConnectionName string `json:"connection_name" bun:"connection_name,type:text,notnull"`
	ConnectionType string `json:"connection_type" bun:"connection_type,type:text,notnull"`
	TenantId       string `json:"tenant_id"`
	BusinessUnitId string `json:"business_unit_id"`
	AppId          string `json:"app_id" bun:"app_id,type:text,notnull"`
	QueueId        string `json:"queue_id" bun:"queue_id,type:text,notnull"`
	Status         bool   `json:"status" bun:"status,notnull"`
}

type ChatConnectionAppRequest struct {
	ConnectionName string `json:"connection_name"`
	ConnectionType string `json:"connection_type"`
	AppId          string `json:"app_id"`
	QueueId        string `json:"queue_id"`
	Status         bool   `json:"status"`
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
	if len(m.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(m.QueueId) < 1 {
		return errors.New("queue id is required")
	}
	return nil
}
