package model

import "github.com/uptrace/bun"

type ConnectionQueue struct {
	*Base
	bun.BaseModel `bun:"table:connection_queue,alias:cq"`
	ConnectionId  string `json:"connection_id" bun:"connection_id,type:text,notnull"`
	QueueId       string `json:"queue_id" bun:"queue_id,type:text,notnull"`
	Status        string `json:"status" bun:"status,notnull"`
}
