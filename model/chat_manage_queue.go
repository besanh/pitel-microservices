package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatManageQueueUser struct {
	*Base
	bun.BaseModel `bun:"table:chat_manage_queue_user,alias:cmqa"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ConnectionId  string `json:"connection_id" bun:"connection_id,type:uuid,default:null"`
	QueueId       string `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	ManageId      string `json:"manage_id" bun:"manage_id,type:uuid,notnull"`
}

type ChatManageQueueUserRequest struct {
	ConnectionId string `json:"connection_id"`
	QueueId      string `json:"queue_id"`
	ManageId     string `json:"manage_id"`
}

func (m *ChatManageQueueUserRequest) Validate() (err error) {
	if len(m.QueueId) < 1 {
		err = errors.New("queue id is required")
	}
	if len(m.ManageId) < 1 {
		err = errors.New("manage id is required")
	}

	return
}
