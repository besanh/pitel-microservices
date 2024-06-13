package model

import (
	"errors"
	"github.com/uptrace/bun"
	"mime/multipart"
)

type ChatScript struct {
	*Base
	bun.BaseModel `bun:"table:chat_script,alias:cst"`
	ScriptName    string             `json:"script_name" bun:"script_name,type:text,notnull"`
	Channel       string             `json:"channel" bun:"channel,type:text,notnull"`
	ConnectionId  string             `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	ConnectionApp *ChatConnectionApp `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	CreatedBy     string             `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy     string             `json:"updated_by" bun:"updated_by,type:uuid,default:null"`
	Status        bool               `json:"status" bun:"status,type:boolean,notnull"`
	ScriptType    string             `json:"script_type" bun:"script_type,type:text,notnull"`
	Content       string             `json:"content,omitempty" bun:"content,type:text"`                 // text script
	FileUrl       string             `json:"file_url,omitempty" bun:"file_url,type:text"`               // file script
	OtherScriptId string             `json:"other_script_id,omitempty" bun:"other_script_id,type:text"` // file script
}

type ChatScriptRequest struct {
	ScriptName    string                `json:"script_name" form:"script_name" binding:"required"`
	Channel       string                `json:"channel" form:"channel" binding:"required"`
	ConnectionId  string                `json:"connection_id" form:"connection_id" binding:"required"`
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
	ScriptName    string             `json:"script_name" bun:"script_name,type:text,notnull"`
	Channel       string             `json:"channel" bun:"channel,type:text,notnull"`
	ConnectionId  string             `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	ConnectionApp *ChatConnectionApp `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	CreatedBy     string             `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy     string             `json:"updated_by" bun:"updated_by,type:uuid,default:null"`
	Status        bool               `json:"status" bun:"status,type:boolean"`
	ScriptType    string             `json:"script_type" bun:"script_type,type:text,notnull"`
	Content       string             `json:"content,omitempty" bun:"content,type:text"`                 // text script
	FileUrl       string             `json:"file_url,omitempty" bun:"file_url,type:text"`               // file script
	OtherScriptId string             `json:"other_script_id,omitempty" bun:"other_script_id,type:uuid"` // file script
}

func (r *ChatScriptRequest) Validate() error {
	if len(r.ScriptName) < 1 {
		return errors.New("script name is required")
	}
	if len(r.ConnectionId) < 1 {
		return errors.New("connection id is required")
	}
	if len(r.ScriptType) < 1 {
		return errors.New("script type is required")
	}
	if len(r.Channel) < 1 {
		return errors.New("channel is required")
	}
	if len(r.Content) < 1 {
		return errors.New("channel is required")
	}

	return nil
}
