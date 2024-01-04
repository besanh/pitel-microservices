package model

import "github.com/uptrace/bun"

type ConnectionQueue struct {
	*Base
	bun.BaseModel `bun:"table:connection_queue,alias:cq"`
	QueueName     string `bun:"queue_name,type:text,notnull"`
}

type QueueAgent struct {
	*Base
	bun.BaseModel `bun:"table:queue_agent,alias:qa"`
	QueueId       string `bun:"queue_id,type:text,notnull"`
	AgentId       string `bun:"agent_id,type:text,notnull"`
}
