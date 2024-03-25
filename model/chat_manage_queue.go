package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatManageQueueAgent struct {
	*Base
	bun.BaseModel `bun:"table:chat_manage_queue_agent,alias:cmqa"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ConnectionId  string `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	QueueId       string `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	ManageId      string `json:"manage_id" bun:"manage_id,type:uuid,notnull"`
}

type ChatManageQueueAgentRequest struct {
	ConnectionId string `json:"connection_id"`
	QueueId      string `json:"queue_id"`
	ManageId     string `json:"manage_id"`
}

func (m *ChatManageQueueAgentRequest) Validate() (err error) {
	if len(m.ConnectionId) < 1 {
		err = errors.New("connection id is required")
	}
	if len(m.QueueId) < 1 {
		err = errors.New("queue id is required")
	}
	if len(m.ManageId) < 1 {
		err = errors.New("manage id is required")
	}

	return
}
