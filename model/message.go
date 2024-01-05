package model

import "time"

type Message struct {
	Id             string         `json:"id"`
	ConversationId string         `json:"conversation_id"`
	TenantId       string         `json:"tenant_id"`
	BusinessUnit   string         `json:"business_unit"`
	UserId         string         `json:"user_id"`
	Username       string         `json:"username"`
	ParentMsgId    string         `json:"parent_msg_id"`
	MsgId          string         `json:"msg_id"`
	MessageType    string         `json:"message_type"`
	EventName      string         `json:"event_name"`
	Direction      string         `json:"direction"`
	AppId          string         `json:"app_id"`
	OaId           string         `json:"oa_id"`          // connection id
	UserIdByApp    string         `json:"user_id_by_app"` // require for zalo
	Uid            string         `json:"uid"`            // id zalo
	UserAppname    string         `json:"user_app_name"`  // username zalo
	Avatar         string         `json:"avatar"`         // avatar zalo
	SupporterId    string         `json:"supporter_id"`   // from crm
	SupporterName  string         `json:"supporter_name"` // from crm
	SendTime       time.Time      `json:"send_time"`
	SendTimestamp  int64          `json:"send_timestamp"`
	Content        string         `json:"content"`
	AgentId        string         `json:"agent_id"`
	AgentName      string         `json:"agent_name"`
	ReadTime       time.Time      `json:"read_time"`
	ReadTimestamp  int64          `json:"read_timestamp"`
	Attachments    []*Attachments `json:"attachments"`
}

type Attachments struct {
	Id                string             `json:"id"`
	MsgId             string             `json:"msg_id"`
	AttachmentType    string             `json:"attachment_type"`
	AttachmentsDetail *AttachmentsDetail `json:"attachments"`
	SendTime          time.Time          `json:"send_time"`
	SendTimestamp     int64              `json:"send_timestamp"`
}

type AttachmentsDetail struct {
	AttachmentMedia *OttPayloadMedia `json:"attachment_media"`
	AttachmentFile  *OttPayloadFile  `json:"attachement_file"`
}
