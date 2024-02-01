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
	AppId   string `json:"app_id"`
	AppName string `json:"app_name"`
	OaId    string `json:"oa_id"`
	AppType string `json:"app_type"`
	Status  string `json:"status"`
}

type ChatConnectionAppFilter struct {
	AppId          string
	ConnectionName string
	ConnectionType string
	QueueId        string
	Status         string
	OaId           string
}

type QueueFilter struct {
	QueueId       []string
	QueueName     string
	ChatRoutingId string
}

type ChatQueueAgentFilter struct {
	QueueId []string
	AgentId []string
	Source  string
}

type ChatRoutingFilter struct {
	RoutingName  string       `json:"routing_name"`
	RoutingAlias string       `json:"routing_alias"`
	Status       sql.NullBool `json:"status"`
}

type ConversationFilter struct {
	AppId          []string
	ConversationId []string
	Username       string
	PhoneNumber    string
	Email          string
}

type AgentAllocationFilter struct {
	ConversationId string
	AgentId        []string
	QueueId        string
	AllocatedTime  int64
}

type MessageFilter struct {
	MessageId           []string        `json:"message_id"`
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
	IsRead              string          `json:"is_read"`
}

type ConnectionQueueFilter struct {
	ConnectionId string `json:"connection_id"`
	QueueId      string `json:"queue_id"`
}
