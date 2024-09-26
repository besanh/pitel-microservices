package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type NotesList struct {
	*Base
	bun.BaseModel  `bun:"table:chat_notes_list"`
	TenantId       string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	Content        string `json:"content" bun:"content,type:text,notnull"`
	ConversationId string `json:"conversation_id" bun:"conversation_id,type:text,notnull"`
	AppId          string `json:"app_id" bun:"app_id,type:text,notnull"`
	OaId           string `json:"oa_id" bun:"oa_id,type:text,notnull"`
}

type ConversationNoteRequest struct {
	Content        string `json:"content"`
	ConversationId string `json:"conversation_id"`
	AppId          string `json:"app_id"`
	OaId           string `json:"oa_id"`
}

func (r *ConversationNoteRequest) Validate() error {
	if len(r.Content) < 1 || len(r.Content) > 1000 {
		return errors.New("content must have number of characters between 1 and 1000")
	}
	if len(r.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(r.OaId) < 1 {
		return errors.New("oa id is required")
	}
	if len(r.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}
	return nil
}
