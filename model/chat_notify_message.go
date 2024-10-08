package model

import (
	"errors"
	"fmt"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type ChatNotifyMessage struct {
	*Base
	bun.BaseModel      `bun:"table:chat_notify_message,alias:cnm"`
	TenantId           string             `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	NotifyType         string             `json:"notify_type" bun:"notify_type,type:text"`
	ConnectionId       string             `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	ConnectionApp      *ChatConnectionApp `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	MessageNotifyAfter int                `json:"message_notify_after" bun:"message_notify_after,type:integer"`
	ReceiverType       string             `json:"receiver_type" bun:"receiver_type,type:text"`
	OaId               string             `json:"oa_id" bun:"oa_id,type:text"`
}

type ChatNotifyMessageRequest struct {
	ConnectionId       string `json:"connection_id"`
	NotifyType         string `json:"notify_type"`
	MessageNotifyAfter int    `json:"message_notify_after"`
	ReceiverType       string `json:"receiver_type"`
}

func (m *ChatNotifyMessageRequest) Validate() error {
	if len(m.ConnectionId) < 1 {
		return errors.New("connection_id is required")
	}
	if !slices.Contains(variables.CHAT_NOTIFY_MESSAGE_TYPE, m.NotifyType) {
		return fmt.Errorf("notify type %s not valid", m.NotifyType)
	}
	if !slices.Contains(variables.CHAT_NOTIFY_MESSAGE_RECEIVER_TYPE, m.ReceiverType) {
		return fmt.Errorf("receiver type %s not valid", m.ReceiverType)
	}
	if m.MessageNotifyAfter < 0 {
		return errors.New("message_notify_after must be greater than zero")
	}

	return nil
}
