package model

import (
	"errors"
	"slices"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
)

type ChatPolicySetting struct {
	*Base
	bun.BaseModel  `bun:"table:chat_policy_setting,alias:cps"`
	TenantId       string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ConnectionType string `json:"connection_type" bun:"connection_type,type:text,notnull"`
	CreatedBy      string `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy      string `json:"updated_by" bun:"updated_by,type:uuid,default:null"`
	ChatWindowTime int    `json:"chat_window_time" bun:"chat_window_time,notnull"`
}

type ChatPolicyConfigRequest struct {
	ConnectionType string `json:"connection_type" binding:"required"`
	ChatWindowTime int    `json:"chat_window_time" binding:"required"`
}

func (c *ChatPolicyConfigRequest) Validate() error {
	if len(c.ConnectionType) < 1 {
		return errors.New("connection_type is required")
	}
	if c.ChatWindowTime < 0 {
		return errors.New("chat_window_time must be greater than zero")
	}
	if !slices.Contains[[]string](variables.CONNECTION_TYPE, c.ConnectionType) {
		return errors.New("connection type " + c.ConnectionType + " is not supported")
	}
	return nil
}
