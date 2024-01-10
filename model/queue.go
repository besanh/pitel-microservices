package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageValueInQueue struct {
	ParentMessageId     *uuid.UUID `json:"parent_message_id"`
	Id                  *uuid.UUID `json:"id"`
	ConversationId      *uuid.UUID `json:"conversation_id"`
	ParentExternalMsgId *uuid.UUID `json:"parent_external_msg_id"`
	ExternalMsgId       *uuid.UUID `json:"externalmsg_id"`
	MessageType         string     `json:"message_type"`
	EventName           string     `json:"event_name"`
	Direction           string     `json:"direction"`
	AppId               *uuid.UUID `json:"app_id"`
	OaId                *uuid.UUID `json:"oa_id"`
	UserIdByApp         *uuid.UUID `json:"user_id_by_app"`
	Uid                 *uuid.UUID `json:"uid"`
	UserAppname         string     `json:"user_app_name"`
	Avatar              string     `json:"avatar"`
	SupporterId         *uuid.UUID `json:"supporter_id"`
	SupporterName       string     `json:"supporter_name"`
	SendTime            time.Time  `json:"send_time"`
	SendTimestamp       int64      `json:"send_timestamp"`
	Content             string     `json:"content"`
	AgentId             *uuid.UUID `json:"agent_id"`
	AgentName           string     `json:"agent_name"`
	ReadTime            time.Time  `json:"read_time"`
	ReadTimestamp       int64      `json:"read_timestamp"`
	Attachments         []any      `json:"attachments"`
}
