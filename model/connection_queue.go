package model

import "github.com/uptrace/bun"

type ConnectionQueue struct {
	*Base
	bun.BaseModel     `bun:"table:chat_connection_queue,alias:ccq"`
	TenantId          string             `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ConnectionId      string             `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	ChatConnectionApp *ChatConnectionApp `json:"chat_connection_app" bun:"rel:has-one,join:connection_id=id"`
	QueueId           string             `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	ChatQueue         *ChatQueue         `json:"chat_queue" bun:"rel:has-one,join:queue_id=id"`
}
