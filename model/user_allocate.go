package model

import (
	"github.com/uptrace/bun"
)

type UserAllocate struct {
	*Base
	bun.BaseModel      `bun:"table:chat_user_allocate,alias:cua"`
	TenantId           string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	AppId              string `json:"app_id" bun:"app_id,type:text,notnull"`
	OaId               string `json:"oa_id" bun:"oa_id,type:text,notnull"`
	ConversationId     string `json:"conversation_id" bun:"conversation_id,type:text,notnull"`
	UserId             string `json:"user_id" bun:"user_id,type:uuid,notnull"`
	QueueId            string `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	AllocatedTimestamp int64  `json:"allocated_timestamp" bun:"allocated_timestamp,notnull"`
	MainAllocate       string `json:"main_allocate" bun:"main_allocate,type:text,notnull"`
	ConnectionId       string `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	Username           string `json:"username" bun:"-"`
}
