package model

import "github.com/uptrace/bun"

// Get agent online on queue to inser
// Table temporary, cronjob 3 hours to delete
type UserAgentAllocation struct {
	*Base
	bun.BaseModel `bun:"table:user_agent_allocation,alias:ua"`
	UserIdByApp   string `bun:"user_id_by_app,type:text,notnull"`
	AgentId       string `bun:"agent_id,type:text,notnull"`
	QueueId       string `bun:"queue_id,type:text,notnull"`
	AllocatedTime int64  `bun:"allocated_time,type:numeric,notnull"`
}
