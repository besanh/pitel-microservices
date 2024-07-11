package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ChatAutoScriptToChatLabel struct {
	bun.BaseModel    `bun:"table:chat_auto_script_to_chat_label,alias:cas_cl"`
	ChatAutoScriptId string          `json:"chat_auto_script_id" bun:"chat_auto_script_id,type:uuid,pk"`
	ChatLabelId      string          `json:"chat_label_id" bun:"chat_label_id,type:uuid,pk"`
	ActionType       string          `json:"action_type" bun:"action_type,type:text,pk"`
	Order            int             `json:"order" bun:"order,notnull"`
	ChatAutoScript   *ChatAutoScript `json:"chat_auto_script" bun:"rel:belongs-to,join:chat_auto_script_id=id"`
	ChatLabel        *ChatLabel      `json:"chat_label" bun:"rel:belongs-to,join:chat_label_id=id"`
	CreatedAt        time.Time       `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt        time.Time       `json:"updated_at" bun:"updated_at,notnull"`
}

type ChatLabelAction struct {
	ChatAutoScriptId string
	ActionType       string
	Order            int
	CreatedAt        time.Time
}
