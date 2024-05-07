package model

import (
	"errors"
	"mime/multipart"
	"time"
)

type Message struct {
	TenantId            string            `json:"tenant_id"`
	ParentMessageId     string            `json:"parent_message_id"`
	Id                  string            `json:"id"`
	ConversationId      string            `json:"conversation_id"`
	ParentExternalMsgId string            `json:"parent_external_msg_id"`
	ExternalMsgId       string            `json:"external_msg_id"`
	MessageType         string            `json:"message_type"`
	EventName           string            `json:"event_name"`
	Direction           string            `json:"direction"`
	AppId               string            `json:"app_id"`
	OaId                string            `json:"oa_id"`            // connection id
	UserIdByApp         string            `json:"user_id_by_app"`   // require for zalo
	ExternalUserId      string            `json:"external_user_id"` // id zalo
	UserAppname         string            `json:"user_app_name"`    // username zalo
	Avatar              string            `json:"avatar"`           // avatar zalo
	SupporterId         string            `json:"supporter_id"`     // from crm
	SupporterName       string            `json:"supporter_name"`   // from crm
	SendTime            time.Time         `json:"send_time"`
	SendTimestamp       int64             `json:"send_timestamp"`
	Content             string            `json:"content"`
	IsRead              string            `json:"is_read"`
	ReadTime            time.Time         `json:"read_time"`
	ReadTimestamp       int64             `json:"read_timestamp"`
	ReadBy              []string          `json:"read_by"`
	Attachments         []*OttAttachments `json:"attachments"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	ShareInfo           *ShareInfo        `json:"share_info"`
}

type AttachmentsDetails struct {
	AttachmentType  string           `json:"attachment_type"`
	AttachmentMedia *OttPayloadMedia `json:"attachment_media"`
	AttachmentFile  *OttPayloadFile  `json:"attachement_file"`
}

type MessageRequest struct {
	EventName      string                `json:"event_name"`
	AppId          string                `json:"app_id"`
	ConversationId string                `json:"conversation_id"`
	Content        string                `json:"content"`
	Attachments    []*AttachmentsDetails `json:"attachments"`
}

type MessageFormRequest struct {
	EventName      string                `form:"event_name" binding:"required"`
	AppId          string                `form:"app_id" binding:"required"`
	ConversationId string                `form:"conversation_id" binding:"required"`
	File           *multipart.FileHeader `form:"file" binding:"required"`
}

type MessageMarkRead struct {
	AppId          string   `json:"app_id"`
	OaId           string   `json:"oa_id"`
	ConversationId string   `json:"conversation_id"`
	MessageIds     []string `json:"message_ids"`
	ReadBy         string   `json:"read_by"`
	ReadAt         string   `json:"read_at"`
	ReadAll        bool     `json:"read_all"`
}

type OaInfoMessage struct {
	ConnectionId string `json:"connection_id"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	Cover        string `json:"cover"`
	CateName     string `json:"cate_name"`
	Code         int64  `json:"code"`
	Message      string `json:"message"`
}

type ShareInfo struct {
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	City        string `json:"city"`
	District    string `json:"district"`
}

type ShareInfoSendToOtt struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	City     string `json:"city"`
	District string `json:"district"`
}

type ReadMessageResponse struct {
	TotalSuccess int               `json:"total_success"`
	TotalFail    int               `json:"total_fail"`
	ListFail     map[string]string `json:"list_fail"`
	ListSuccess  map[string]string `json:"list_success"`
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

func (m *MessageMarkRead) ValidateMarkRead() error {
	if len(m.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}
	if !m.ReadAll {
		if len(m.MessageIds) < 1 {
			return errors.New("message ids is required")
		}
	}

	return nil
}
