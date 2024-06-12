package model

import (
	"errors"
	"github.com/uptrace/bun"
	"mime/multipart"
)

type ChatCommand struct {
	*Base
	bun.BaseModel `bun:"table:chat_command,alias:cc"`
	Keyword       string `json:"keyword" bun:"keyword,type:text,notnull"`
	Theme         string `json:"theme" bun:"theme,type:text,notnull"`
	PageId        string `json:"page_id" bun:"page_id,type:uuid,notnull"`
	Channel       string `json:"channel" bun:"channel,type:text,notnull"`
	Content       string `json:"content" bun:"content,type:text,notnull"`
	CreatorId     string `json:"creator_id" bun:"creator_id,type:uuid,notnull"`
	ImageUrl      string `json:"image_url,omitempty" bun:"image_url,type:text"`
}

type ChatPersonalization struct {
	*Base
	bun.BaseModel `bun:"table:chat_personalization,alias:cp"`
	Value         string `json:"value" json:"value,type:text,notnull"`
}

type ChatCommandRequest struct {
	Keyword string                `form:"keyword" binding:"required"`
	Theme   string                `form:"theme" binding:"required"`
	PageId  string                `form:"page_id" binding:"required"`
	Channel string                `form:"channel" binding:"required"`
	Content string                `form:"content" binding:"required"`
	File    *multipart.FileHeader `form:"file"`
}

type ChatCommandView struct {
	*ChatCommand
	ConnectionName string `json:"connection_name" bun:"connection_name,type=text"`
}

func (r *ChatCommandRequest) Validate() error {
	if len(r.Keyword) < 1 {
		return errors.New("keyword is required")
	}
	if len(r.Theme) < 1 {
		return errors.New("theme is required")
	}
	if len(r.PageId) < 1 {
		return errors.New("page id is required")
	}
	if len(r.Channel) < 1 {
		return errors.New("channel is required")
	}
	if len(r.Content) < 1 {
		return errors.New("channel is required")
	}

	return nil
}
