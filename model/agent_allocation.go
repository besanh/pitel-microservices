package model

import (
	"github.com/uptrace/bun"
)

// Get agent online on queue to insert
// Table temporary, cronjob 3 hours to delete
type AgentAllocation struct {
	*Base
	bun.BaseModel      `bun:"table:agent_allocation,alias:aa"`
	TenantId           string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	AppId              string `json:"app_id" bun:"app_id,type:text"`
	ConversationId     string `json:"conversation_id" bun:"conversation_id,type:text,notnull"`
	AgentId            string `json:"agent_id" bun:"agent_id,type:text,notnull"`
	QueueId            string `json:"queue_id" bun:"queue_id,type:text,notnull"`
	AllocatedTimestamp int64  `json:"allocated_timestamp" bun:"allocated_timestamp,notnull"`
	MainAllocate       string `json:"main_allocate" bun:"main_allocate,type:string,notnull"`
	Source             string `json:"source" bun:"source,type:text,notnull"`
}
