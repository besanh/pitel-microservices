package model

import (
	"database/sql"
	"encoding/json"
)

type AuthSourceFilter struct {
	Source string
	Status sql.NullBool
}

type AppFilter struct {
	AppName string
	Status  sql.NullBool
}

type ChatConnectionAppFilter struct {
	ConnectionName string
	ConnectionType string
	Status         sql.NullBool
}

type QueueFilter struct {
	AppId     string
	QueueName string
}

type ChatQueueAgentFilter struct {
	QueueId string
	AgentId string
	Source  string
}

type ChatRoutingFilter struct {
	RoutingName string
	Status      sql.NullBool
}

type ConversationFilter struct {
	AppId          []string
	ConversationId []string
	Username       []string
	PhoneNumber    []string
	Email          []string
}

type AgentAllocationFilter struct {
	ConversationId string
	AgentId        []string
	QueueId        string
	AllocatedTime  int64
}

type MessageFilter struct {
	ParentMessageId     string          `json:"parent_message_id"`
	ConversationId      string          `json:"conversation_id"`
	ParentExternalMsgId string          `json:"parent_external_msg_id"`
	MessageType         string          `json:"message_type"`
	EventName           string          `json:"event_name"`
	Direction           string          `json:"direction"`
	AppId               string          `json:"app_id"`
	OaId                string          `json:"oa_id"`
	UserIdByApp         string          `json:"user_id_by_app"`
	ExternalUserId      string          `json:"external_user_id"`
	UserAppname         string          `json:"user_app_name"`
	SupporterId         string          `json:"supporter_id"`
	SupporterName       string          `json:"supporter_name"`
	SendTime            string          `json:"send_time"`
	Content             string          `json:"content"`
	ReadTime            string          `json:"read_time"`
	ReadBy              json.RawMessage `json:"read_by"`
	IsRead              bool            `json:"is_read"`
}
