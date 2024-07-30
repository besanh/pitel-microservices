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

type ChatAppFilter struct {
	AppId   string `json:"app_id"`
	AppName string `json:"app_name"`
	OaId    string `json:"oa_id"`
	AppType string `json:"app_type"`
	Status  string `json:"status"`
}

type ChatConnectionAppFilter struct {
	AppId             string `json:"app_id"`
	TenantId          string `json:"tenant_id"`
	ConnectionName    string `json:"connection_name"`
	ConnectionType    string `json:"connection_type"`
	QueueId           string `json:"queue_id"`
	ConnectionQueueId string `json:"connection_queue_id"`
	Status            string `json:"status"`
	OaId              string `json:"oa_id"`
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
	RoutingName  string       `json:"routing_name"`
	RoutingAlias string       `json:"routing_alias"`
	Status       sql.NullBool `json:"status"`
}

type ConversationFilter struct {
	AppId                  []string     `json:"app_id"`
	TenantId               string       `json:"tenant_id"`
	ConversationId         []string     `json:"conversation_id"`
	ExternalConversationId []string     `json:"external_conversation_id"`
	Username               string       `json:"username"`
	PhoneNumber            string       `json:"phone_number"`
	Email                  string       `json:"email"`
	Insensitive            string       `json:"insensitive"`
	IsDone                 sql.NullBool `json:"is_done"`
	Major                  sql.NullBool `json:"major"`
	Following              sql.NullBool `json:"following"`
}

type AllocateUserFilter struct {
	TenantId               string   `json:"tenant_id"`
	AppId                  string   `json:"app_id"`
	OaId                   string   `json:"oa_id"`
	ConversationId         string   `json:"conversation_id"`
	ExternalConversationId string   `json:"external_conversation_id"`
	UserId                 []string `json:"user_id"`
	QueueId                []string `json:"queue_id"`
	AllocatedTime          int64    `json:"allocated_time"`
	MainAllocate           string   `json:"main_allocate"`
}

type MessageFilter struct {
	TenantId            string          `json:"tenant_id"`
	MessageId           []string        `json:"message_id"`
	ParentMessageId     string          `json:"parent_message_id"`
	ConversationId      string          `json:"conversation_id"`
	ParentExternalMsgId string          `json:"parent_external_msg_id"`
	ExternalMessageId   string          `json:"external_message_id"`
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
	TenantId          string `json:"tenant_id"`
	ConnectionId      string `json:"connection_id"`
	QueueId           string `json:"queue_id"`
	ConnectionQueueId string `json:"connection_queue_id"`
}

type ShareInfoFormFilter struct {
	TenantId     string `json:"tenant_id"`
	ShareType    string `json:"share_type"`
	AppId        string `json:"app_id"`
	OaId         string `json:"oa_id"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	ConnectionId string `json:"connection_id"`
}

type FacebookPageFilter struct {
	OaId   string `json:"oa_id"`
	OaName string `json:"oa_name"`
}

type ChatManageQueueUserFilter struct {
	TenantId     string `json:"tenant_id"`
	ConnectionId string `json:"connection_id"`
	QueueId      string `json:"queue_id"`
	UserId       string `json:"user_id"`
}

type UserInQueueFilter struct {
	AppId            string `json:"app_id"`
	OaId             string `json:"oa_id"`
	ConversationId   string `json:"conversation_id"`
	ConversationType string `json:"conversation_type"`
	Status           string `json:"status"`
}

type ChatEmailFilter struct {
	TenantId string       `json:"tenant_id"`
	OaId     string       `json:"oa_id"`
	Status   sql.NullBool `json:"status"`
}

type ChatMsgSampleFilter struct {
	TenantId     string `json:"tenant_id"`
	ConnectionId string `json:"connection_id"`
	OaId         string `json:"oa_id"`
	Channel      string `json:"channel"`
	Keyword      string `json:"keyword"`
}

type ChatScriptFilter struct {
	TenantId   string       `json:"tenant_id"`
	Channel    string       `json:"channel"`
	Status     sql.NullBool `json:"status"`
	ScriptName string       `json:"script_name"`
}

type ChatLabelFilter struct {
	TenantId         string       `json:"tenant_id"`
	AppId            string       `json:"app_id"`
	OaId             string       `json:"oa_id"`
	LabelName        string       `json:"label_name"`
	LabelType        string       `json:"label_type"`
	LabelColor       string       `json:"label_color"`
	LabelStatus      sql.NullBool `json:"label_status"`
	ExternalLabelId  string       `json:"external_label_id"`
	IsSearchExactly  sql.NullBool `json:"is_search_exactly"`
	IsSearchMultiple sql.NullBool `json:"is_search_multiple"`
	LabelIds         []string     `json:"label_ids"`
}

type ChatAutoScriptFilter struct {
	TenantId     string       `json:"tenant_id"`
	Channel      string       `json:"channel"`
	ScriptName   string       `json:"script_name"`
	Status       sql.NullBool `json:"status"`
	OaId         string       `json:"oa_id"`
	TriggerEvent string       `json:"trigger_event"`
}

type ChatPolicyFilter struct {
	TenantId       string   `json:"tenant_id"`
	ConnectionType string   `json:"connection_type"`
	ExcludedIds    []string `json:"excluded_ids"`
}

type ChatIntegrateSystemFilter struct {
	SystemName      string       `json:"system_name"`
	VendorName      string       `json:"vendor_name"`
	Status          sql.NullBool `json:"status"`
	SystemId        string       `json:"system_id"`
	TenantDefaultId string       `json:"tenant_default_id"`
	ServerName      string       `json:"server_name"`
}

type ChatVendorFilter struct {
	VendorName string       `json:"vendor_name"`
	VendorType string       `json:"vendor_type"`
	Status     sql.NullBool `json:"status"`
}

type ChatRoleFilter struct {
	RoleName string       `json:"role_name"`
	Status   sql.NullBool `json:"status"`
}

type ChatUserFilter struct {
	Username string       `json:"username"`
	Level    string       `json:"level"`
	Status   sql.NullBool `json:"status"`
	Fullname string       `json:"fullname"`
	RoleId   string       `json:"role_id"`
}

type ChatTenantFilter struct {
	TenantName string       `json:"tenant_name"`
	Status     sql.NullBool `json:"status"`
}

type ChatAppIntegrateSystemFilter struct {
	ChatAppId             string `json:"chat_app_id"`
	ChatIntegrateSystemId string `json:"chat_integrate_system_id"`
}
