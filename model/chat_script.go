package model

import (
	"errors"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
	"mime/multipart"
)

type ChatScript struct {
	*Base
	bun.BaseModel `bun:"table:chat_script,alias:cst"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ScriptName    string `json:"script_name" bun:"script_name,type:text,notnull"`
	Channel       string `json:"channel" bun:"channel,type:text,notnull"`
	CreatedBy     string `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy     string `json:"updated_by" bun:"updated_by,type:uuid,default:null"`
	Status        bool   `json:"status" bun:"status,type:boolean,notnull"`
	ScriptType    string `json:"script_type" bun:"script_type,type:text,notnull"`
	Content       string `json:"content" bun:"content,type:text"`   // text script
	FileUrl       string `json:"file_url" bun:"file_url,type:text"` // file script
	OtherScriptId string `json:"other_script_id" bun:"other_script_id,type:text"`
}

type ChatScriptRequest struct {
	ScriptName    string                `json:"script_name" form:"script_name" binding:"required"`
	Channel       string                `json:"channel" form:"channel" binding:"required"`
	Status        string                `json:"status" form:"status" binding:"required"`
	ScriptType    string                `json:"script_type" form:"script_type" binding:"required"`
	Content       string                `form:"content"`
	File          *multipart.FileHeader `form:"file"`
	OtherScriptId string                `form:"other_script_id"`
}

type ChatScriptStatusRequest struct {
	Status string `json:"status" form:"status" binding:"required"`
}

type ChatScriptView struct {
	*Base
	bun.BaseModel `bun:"table:chat_script,alias:cst"`
	TenantId      string `json:"tenant_id" bun:"tenant_id"`
	ScriptName    string `json:"script_name" bun:"script_name"`
	Channel       string `json:"channel" bun:"channel"`
	CreatedBy     string `json:"created_by" bun:"created_by"`
	UpdatedBy     string `json:"updated_by" bun:"updated_by"`
	Status        bool   `json:"status" bun:"status"`
	ScriptType    string `json:"script_type" bun:"script_type"`
	Content       string `json:"content" bun:"content"`   // text script
	FileUrl       string `json:"file_url" bun:"file_url"` // file script
	OtherScriptId string `json:"other_script_id" bun:"other_script_id"`
}

func (r *ChatScriptRequest) Validate() error {
	if len(r.ScriptName) < 1 {
		return errors.New("script name is required")
	}
	if len(r.ScriptType) < 1 {
		return errors.New("script type is required")
	}
	if len(r.Channel) < 1 {
		return errors.New("channel is required")
	}
	if !slices.Contains[[]string](variables.CHAT_SCRIPT_TYPE, r.ScriptType) {
		return errors.New("script type " + r.ScriptType + " is not supported")
	}

	return nil
}
