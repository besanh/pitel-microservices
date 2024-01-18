package model

import (
	"encoding/json"
	"errors"
	"time"
)

type Message struct {
	ParentMessageId     string          `json:"parent_message_id"`
	Id                  string          `json:"id"`
	ConversationId      string          `json:"conversation_id"`
	ParentExternalMsgId string          `json:"parent_external_msg_id"`
	ExternalMsgId       string          `json:"external_msg_id"`
	MessageType         string          `json:"message_type"`
	EventName           string          `json:"event_name"`
	Direction           string          `json:"direction"`
	AppId               string          `json:"app_id"`
	OaId                string          `json:"oa_id"`            // connection id
	UserIdByApp         string          `json:"user_id_by_app"`   // require for zalo
	ExternalUserId      string          `json:"external_user_id"` // id zalo
	UserAppname         string          `json:"user_app_name"`    // username zalo
	Avatar              string          `json:"avatar"`           // avatar zalo
	SupporterId         string          `json:"supporter_id"`     // from crm
	SupporterName       string          `json:"supporter_name"`   // from crm
	SendTime            time.Time       `json:"send_time"`
	SendTimestamp       int64           `json:"send_timestamp"`
	Content             string          `json:"content"`
	ReadTime            time.Time       `json:"read_time"`
	ReadTimestamp       int64           `json:"read_timestamp"`
	ReadBy              json.RawMessage `json:"read_by"`
	Attachments         []*Attachments  `json:"attachments"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
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

type MessageRequest struct {
	AppId          string         `json:"app_id"`
	ConversationId string         `json:"conversation_id"`
	Content        string         `json:"content"`
	Attachments    []*Attachments `json:"attachments"`
}

func (m *MessageRequest) Validate() error {
	if len(m.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}
	if len(m.Content) < 1 {
		return errors.New("content is required")
	}
	return nil
}
