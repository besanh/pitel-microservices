package model

import "github.com/uptrace/bun"

type ConnectionQueue struct {
	*Base
	bun.BaseModel     `bun:"table:connection_queue,alias:cq"`
	ConnectionId      string               `json:"connection_id" bun:"connection_id,type:text,notnull"`
	ChatConnectionApp []*ChatConnectionApp `json:"chat_connection_app" bun:"rel:has-many,join:connection_id=id"`
	QueueId           string               `json:"queue_id" bun:"queue_id,type:text,notnull"`
	ChatQueue         *ChatQueue           `json:"chat_queue" bun:"rel:has-one,join:queue_id=id"`
	Status            string               `json:"status" bun:"status,notnull"`
}
