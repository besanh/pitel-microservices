package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatQueueAgent struct {
	*Base
	bun.BaseModel  `bun:"table:chat_queue_agent,alias:qa"`
	TenantId       string     `json:"tenant_id" bun:"tenant_id,type:text,notnull"`
	BusinessUnitId string     `json:"business_unit_id" bun:"business_unit_id,type:text,notnull"`
	QueueId        string     `json:"queue_id" bun:"queue_id,type:text,notnull"`
	ChatQueue      *ChatQueue `json:"chat_queue" bun:"rel:belongs-to,join:queue_id=id"`
	AgentId        string     `json:"agent_id" bun:"agent_id,type:text,notnull"`
	Source         string     `json:"source" bun:"source,type:text,notnull"`
}

type ChatQueueAgentRequest struct {
	QueueId string `json:"queue_id"`
	AgentId string `json:"agent_id"`
	Source  string `json:"source"`
}

func (m *ChatQueueAgentRequest) Validate() error {
	if len(m.QueueId) < 1 {
		return errors.New("queue id is required")
	}
	if len(m.AgentId) < 1 {
		return errors.New("agent id is required")
	}
	if len(m.Source) < 1 {
		return errors.New("source is required")
	}
	return nil
}
