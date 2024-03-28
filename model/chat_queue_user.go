package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatQueueUser struct {
	*Base
	bun.BaseModel `bun:"table:chat_queue_user,alias:cqa"`
	TenantId      string     `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	QueueId       string     `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	ChatQueue     *ChatQueue `json:"chat_queue" bun:"rel:belongs-to,join:queue_id=id"`
	UserId        string     `json:"user_id" bun:"user_id,type:text,notnull"`
	Source        string     `json:"source" bun:"source,type:text,notnull"`
	Fullname      string     `json:"fullname" bun:"-"`
}

type ChatQueueUserView struct {
	*Base
	bun.BaseModel `bun:"table:chat_queue_user_view,alias:cqa"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	QueueId       string `json:"queue_id" bun:"queue_id,type:uuid,notnull"`
	UserId        string `json:"user_id" bun:"user_id,type:text,notnull"`
}

type ChatQueueUserRequest struct {
	QueueId string   `json:"queue_id"`
	UserId  []string `json:"user_id"`
	Source  string   `json:"source"`
}

type ChatQueueUserUpdateResponse struct {
	TotalSuccess int      `json:"total_success"`
	TotalFail    int      `json:"total_fail"`
	ListFail     []string `json:"list_fail"`
	ListSuccess  []string `json:"list_success"`
}

func (m *ChatQueueUserRequest) Validate() error {
	if len(m.QueueId) < 1 {
		return errors.New("queue id is required")
	}
	if len(m.UserId) < 1 {
		return errors.New("user id is required")
	}
	if len(m.Source) < 1 {
		return errors.New("source is required")
	}
	return nil
}
