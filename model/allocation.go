package model

import "github.com/uptrace/bun"

// Get agent online on queue to insert
// Table temporary, cronjob 3 hours to delete
type AgentAllocation struct {
	*Base
	bun.BaseModel `bun:"table:agent_allocation,alias:ua"`
	UserIdByApp   string `bun:"user_id_by_app,type:text,notnull"`
	AgentId       string `bun:"agent_id,type:text,notnull"`
	QueueId       string `bun:"queue_id,type:text,notnull"`
	AllocatedTime int64  `bun:"allocated_time,type:numeric,notnull"`
}
