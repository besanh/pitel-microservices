package model

import (
	"database/sql"
	"encoding/json"
)

type AuthSourceFilter struct {
	TenantId string `json:"tenant_id"`
	Source   string
	Status   sql.NullBool
}

type AppFilter struct {
	AppId      string `json:"app_id"`
	AppName    string `json:"app_name"`
	OaId       string `json:"oa_id"`
	AppType    string `json:"app_type"`
	Status     string `json:"status"`
	DefaultApp string `json:"default_app"`
}

type ChatConnectionAppFilter struct {
	AppId          string
	TenantId       string `json:"tenant_id"`
	ConnectionName string
	ConnectionType string
	QueueId        string
	Status         string
	OaId           string
}

type QueueFilter struct {
	TenantId      string `json:"tenant_id"`
	QueueId       []string
	QueueName     string
	ChatRoutingId string
}

type ChatQueueUserFilter struct {
	TenantId string   `json:"tenant_id"`
	QueueId  []string `json:"queue_id"`
	UserId   []string `json:"user_id"`
	Source   string   `json:"source"`
}

type ChatRoutingFilter struct {
	TenantId     string       `json:"tenant_id"`
	RoutingName  string       `json:"routing_name"`
	RoutingAlias string       `json:"routing_alias"`
	Status       sql.NullBool `json:"status"`
}

type ConversationFilter struct {
	AppId          []string     `json:"app_id"`
	TenantId       string       `json:"tenant_id"`
	ConversationId []string     `json:"conversation_id"`
	Username       string       `json:"username"`
	PhoneNumber    string       `json:"phone_number"`
	Email          string       `json:"email"`
	Insensitive    string       `json:"insensitive"`
	IsDone         sql.NullBool `json:"is_done"`
}

type UserAllocateFilter struct {
	TenantId       string   `json:"tenant_id"`
	AppId          string   `json:"app_id"`
	OaId           string   `json:"oa_id"`
	ConversationId string   `json:"conversation_id"`
	UserId         []string `json:"user_id"`
	QueueId        string   `json:"queue_id"`
	AllocatedTime  int64    `json:"allocated_time"`
	MainAllocate   string   `json:"main_allocate"`
}

type MessageFilter struct {
	TenantId            string          `json:"tenant_id"`
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
	EventNameExlucde    []string        `json:"event_name_exclude"`
}

type ConnectionQueueFilter struct {
	TenantId     string `json:"tenant_id"`
	ConnectionId string `json:"connection_id"`
	QueueId      string `json:"queue_id"`
}

type ShareInfoFormFilter struct {
	TenantId  string `json:"tenant_id"`
	ShareType string `json:"share_type"`
	AppId     string `json:"app_id"`
	OaId      string `json:"oa_id"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
}

type FacebookPageFilter struct {
	OaId   string `json:"oa_id"`
	OaName string `json:"oa_name"`
}

type ChatManageQueueUserFilter struct {
	ConnectionId string `json:"connection_id"`
	QueueId      string `json:"queue_id"`
	ManageId     string `json:"manage_id"`
}

type UserInQueueFilter struct {
	AppId            string `json:"app_id"`
	OaId             string `json:"oa_id"`
	ConversationId   string `json:"conversation_id"`
	ConversationType string `json:"conversation_type"`
	Status           string `json:"status"`
}
