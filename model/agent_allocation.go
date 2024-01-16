package model

import (
	"github.com/uptrace/bun"
)

// Get agent online on queue to insert
// Table temporary, cronjob 3 hours to delete
type AgentAllocation struct {
	*Base
	bun.BaseModel      `bun:"table:agent_allocation,alias:aa"`
	UserIdByApp        string `json:"user_id_by_app" bun:"user_id_by_app,type:text,notnull"`
	AgentId            string `json:"agent_id" bun:"agent_id,type:text,notnull"`
	QueueId            string `json:"queue_id" bun:"queue_id,type:text,notnull"`
	AllocatedTimestamp int64  `json:"allocated_timestamp" bun:"allocated_timestamp,notnull"`
}
