package model

import (
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type MediaType string

const (
	Media        MediaType = "media"
	MediaImage   MediaType = "image"
	MediaAudio   MediaType = "audio"
	MediaVideo   MediaType = "video"
	MediaLink    MediaType = "link"
	MediaSticker MediaType = "sticker"
	MediaGif     MediaType = "gif"
	MediaFile    MediaType = "file"
)

type ConversationMedia struct {
	*Base
	bun.BaseModel          `bun:"table:conversation_media,alias:com"`
	TenantId               string    `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ConversationId         string    `json:"conversation_id" bun:"conversation_id,type:text,notnull"`
	ExternalConversationId string    `json:"external_conversation_id" bun:"external_conversation_id,type:text,notnull"`
	ConversationType       string    `json:"conversation_type" bun:"conversation_type,type:text,notnull"`
	MessageId              string    `json:"message_id" bun:"message_id,type:text,notnull"`
	MediaType              string    `json:"media_type" bun:"media_type,type:text,notnull"`
	MediaHeader            string    `json:"media_header" bun:"media_header,type:text"`
	MediaUrl               string    `json:"media_url" bun:"media_url,type:text"`
	MediaSize              int64     `json:"media_size" bun:"media_size,type:integer"`
	SendTimestamp          time.Time `json:"send_timestamp" bun:"send_timestamp,notnull"`
}

type ConversationMediaRequest struct {
	TenantId               string    `json:"tenant_id"`
	ConversationId         string    `json:"conversation_id"`
	ExternalConversationId string    `json:"external_conversation_id"`
	ConversationType       string    `json:"conversation_type"`
	MessageId              string    `json:"message_id"`
	MediaType              string    `json:"media_type"`
	MediaHeader            string    `json:"media_header"`
	MediaUrl               string    `json:"media_url"`
	MediaSize              string    `json:"media_size"`
	SendTimestamp          time.Time `json:"send_timestamp"`
}

func (m *ConversationMediaRequest) Validate() error {
	if len(m.TenantId) < 1 {
		return errors.New("tenant id is required")
	}
	if len(m.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}
	if len(m.ExternalConversationId) < 1 {
		return errors.New("external conversation id is required")
	}
	if len(m.ConversationType) < 1 {
		return errors.New("conversation type is required")
	}
	if len(m.MessageId) < 1 {
		return errors.New("message id is required")
	}
	if len(m.MediaType) < 1 {
		return errors.New("media type is required")
	}
	switch MediaType(m.MediaType) {
	case MediaImage, MediaAudio, MediaVideo, MediaSticker, MediaGif, MediaLink, MediaFile:
		if len(m.MediaUrl) < 1 {
			return errors.New("media url is required")
		}
	default:
		return errors.New("unsupported media type")
	}

	return nil
}

func (m *ConversationMediaFilter) Validate() error {
	if len(m.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}
	return nil
}
