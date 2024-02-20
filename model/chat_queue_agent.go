package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatQueueAgent struct {
	*Base
	bun.BaseModel `bun:"table:chat_queue_agent,alias:qa"`
	TenantId      string     `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	QueueId       string     `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	ChatQueue     *ChatQueue `json:"chat_queue" bun:"rel:belongs-to,join:queue_id=id"`
	AgentId       string     `json:"agent_id" bun:"agent_id,type:text,notnull"`
	Source        string     `json:"source" bun:"source,type:text,notnull"`
	Fullname      string     `json:"fullname" bun:"-"`
}

type ChatQueueAgentRequest struct {
	QueueId string   `json:"queue_id"`
	AgentId []string `json:"agent_id"`
	Source  string   `json:"source"`
}

type ChatQueueAgentUpdateResponse struct {
	TotalSuccess int      `json:"total_success"`
	TotalFail    int      `json:"total_fail"`
	ListFail     []string `json:"list_fail"`
	ListSuccess  []string `json:"list_success"`
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
