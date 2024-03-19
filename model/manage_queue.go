package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatManageQueueAgent struct {
	*Base
	bun.BaseModel `bun:"table:chat_manage_queue_agent,alias:cmqa"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	QueueId       string `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	AgentId       string `json:"agent_id" bun:"agent_id,type:uuid,notnull"`
}

type ManageQueueAgentRequest struct {
	QueueId string `json:"queue_id"`
	AgentId string `json:"agent_id"`
}

func (m *ManageQueueAgentRequest) Validate() (err error) {
	if len(m.QueueId) < 1 {
		err = errors.New("queue id is required")
	}
	if len(m.AgentId) < 1 {
		err = errors.New("agent id is required")
	}

	return
}
