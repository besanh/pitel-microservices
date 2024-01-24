package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatQueue struct {
	*Base
	bun.BaseModel `bun:"table:chat_queue,alias:cq"`
	QueueName     string       `json:"queue_name" bun:"queue_name,type:text,notnull"`
	Description   string       `json:"description" bun:"description,type:text"`
	ChatRoutingId string       `json:"chat_routing_id" bun:"chat_routing_id,type:text,notnull"`
	ChatRouting   *ChatRouting `json:"chat_routing" bun:"rel:has-one,join:chat_routing_id=id"`
	Status        bool         `json:"status" bun:"status,notnull"`
}

type ChatQueueRequest struct {
	QueueName     string `json:"queue_name"`
	Description   string `json:"description"`
	ChatRoutingId string `json:"chat_routing_id"`
	Status        bool   `json:"status"`
}

func (m *ChatQueueRequest) Validate() error {
	if len(m.QueueName) < 1 {
		return errors.New("queue name is required")
	}
	if len(m.ChatRoutingId) < 1 {
		return errors.New("chat routing id is required")
	}
	return nil
}
