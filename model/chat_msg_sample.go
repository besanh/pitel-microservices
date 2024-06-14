package model

import (
	"errors"
	"github.com/uptrace/bun"
	"mime/multipart"
)

type ChatMsgSample struct {
	*Base
	bun.BaseModel `bun:"table:chat_message_sample,alias:cms"`
	Keyword       string             `json:"keyword" bun:"keyword,type:text,notnull"`
	Theme         string             `json:"theme" bun:"theme,type:text,notnull"`
	ConnectionId  string             `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	ConnectionApp *ChatConnectionApp `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	Channel       string             `json:"channel" bun:"channel,type:text,notnull"`
	Content       string             `json:"content" bun:"content,type:text,notnull"`
	CreatedBy     string             `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy     string             `json:"updated_by" bun:"updated_by,type:uuid,default:null"`
	ImageUrl      string             `json:"image_url" bun:"image_url,type:text"`
}

type ChatMsgSampleRequest struct {
	Keyword      string                `form:"keyword" binding:"required"`
	Theme        string                `form:"theme" binding:"required"`
	ConnectionId string                `json:"connection_id" form:"connection_id" binding:"required"`
	Channel      string                `form:"channel" binding:"required"`
	Content      string                `form:"content" binding:"required"`
	File         *multipart.FileHeader `form:"file"`
}

type ChatMsgSampleView struct {
	*Base
	bun.BaseModel `bun:"table:chat_message_sample,alias:cms"`
	Keyword       string             `json:"keyword" bun:"keyword"`
	Theme         string             `json:"theme" bun:"theme"`
	ConnectionId  string             `json:"connection_id" bun:"connection_id"`
	ConnectionApp *ChatConnectionApp `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	Channel       string             `json:"channel" bun:"channel"`
	Content       string             `json:"content" bun:"content"`
	CreatedBy     string             `json:"created_by" bun:"created_by"`
	UpdatedBy     string             `json:"updated_by" bun:"updated_by"`
	ImageUrl      string             `json:"image_url" bun:"image_url"`
}

func (r *ChatMsgSampleRequest) Validate() error {
	if len(r.Keyword) < 1 {
		return errors.New("keyword is required")
	}
	if len(r.Theme) < 1 {
		return errors.New("theme is required")
	}
	if len(r.ConnectionId) < 1 {
		return errors.New("connection id is required")
	}
	if len(r.Channel) < 1 {
		return errors.New("channel is required")
	}
	if len(r.Content) < 1 {
		return errors.New("channel is required")
	}

	return nil
}
