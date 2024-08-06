package model

import (
	"errors"
	"mime/multipart"
	"time"
)

type Message struct {
	TenantId               string            `json:"tenant_id"`
	ParentMessageId        string            `json:"parent_message_id"`
	MessageId              string            `json:"message_id"`
	ConversationId         string            `json:"conversation_id"` // uuid of conversation
	ParentExternalMsgId    string            `json:"parent_external_msg_id"`
	ExternalConversationId string            `json:"external_conversation_id"` // app_id+oa_id+external_user_id
	ExternalMsgId          string            `json:"external_msg_id"`
	MessageType            string            `json:"message_type"`
	EventName              string            `json:"event_name"`
	Direction              string            `json:"direction"`
	AppId                  string            `json:"app_id"`
	OaId                   string            `json:"oa_id"`            // connection id
	UserIdByApp            string            `json:"user_id_by_app"`   // require for zalo
	ExternalUserId         string            `json:"external_user_id"` // id zalo
	UserAppname            string            `json:"user_app_name"`    // username zalo
	Avatar                 string            `json:"avatar"`           // avatar zalo
	SupporterId            string            `json:"supporter_id"`     // from crm
	SupporterName          string            `json:"supporter_name"`   // from crm
	SendTime               time.Time         `json:"send_time"`
	SendTimestamp          int64             `json:"send_timestamp"`
	Content                string            `json:"content"`
	IsRead                 string            `json:"is_read"`
	ReadTime               time.Time         `json:"read_time"`
	ReadTimestamp          int64             `json:"read_timestamp"`
	ReadBy                 []string          `json:"read_by"`
	Attachments            []*OttAttachments `json:"attachments"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	ShareInfo              *ShareInfo        `json:"share_info"`
	IsEcho                 bool              `json:"is_echo"`
}

type AttachmentsDetails struct {
	AttachmentType  string           `json:"attachment_type"`
	AttachmentMedia *OttPayloadMedia `json:"attachment_media"`
	AttachmentFile  *OttPayloadFile  `json:"attachement_file"`
}

type MessageAttachmentsDetails struct {
	AttachmentType string          `json:"att_type"`
	Payload        OttPayloadMedia `json:"payload"`
	MessageId      string          `json:"message_id"`
}

type MessageRequest struct {
	EventName      string                `json:"event_name"`
	AppId          string                `json:"app_id"`
	OaId           string                `json:"oa_id"`
	ConversationId string                `json:"conversation_id"`
	Content        string                `json:"content"`
	Attachments    []*AttachmentsDetails `json:"attachments"`
	Url            string                `json:"url"`
}

type MessageFormRequest struct {
	EventName      string                `form:"event_name" binding:"required"`
	AppId          string                `form:"app_id" binding:"required"`
	OaId           string                `form:"oa_id" binding:"required"`
	ConversationId string                `form:"conversation_id" binding:"required"`
	File           *multipart.FileHeader `form:"file"`
	Url            string                `form:"url"`
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
	ConnectionId        string  `json:"connection_id"`
	Name                string  `json:"name"`
	Avatar              string  `json:"avatar"`
	Cover               string  `json:"cover"`
	CateName            string  `json:"cate_name"`
	Code                int64   `json:"code"`
	Message             string  `json:"message"`
	AccessToken         string  `json:"access_token"`
	Expire              int64   `json:"expire"`
	TokenCreatedAt      string  `json:"token_created_at"`
	TokenExpiresIn      float64 `json:"token_expires_in"`
	TokenTimeRemainning float64 `json:"token_time_remaining"`
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

func (m *MessageFormRequest) ValidateMessageForm() error {
	if len(m.EventName) < 1 {
		return errors.New("event name is required")
	}
	if len(m.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(m.OaId) < 1 {
		return errors.New("oa id is required")
	}
	if len(m.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}

	if len(m.Url) < 1 && m.File == nil {
		return errors.New("url or file is required")
	}
	return nil
}
