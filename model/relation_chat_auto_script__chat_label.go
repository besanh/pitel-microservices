package model

import (
	"github.com/uptrace/bun"
	"time"
)

type ChatAutoScriptToChatLabel struct {
	bun.BaseModel    `bun:"table:chat_auto_script_to_chat_label,alias:cas_cl"`
	ChatAutoScriptId string          `bun:"chat_auto_script_id,type:uuid,pk"`
	ChatLabelId      string          `bun:"chat_label_id,type:uuid,pk"`
	ActionType       string          `bun:"action_type,type:text,pk"`
	Order            int             `bun:"order,notnull"`
	ChatAutoScript   *ChatAutoScript `bun:"rel:belongs-to,join:chat_auto_script_id=id"`
	ChatLabel        *ChatLabel      `bun:"rel:belongs-to,join:chat_label_id=id"`
	CreatedAt        time.Time       `bun:"created_at,notnull"`
	UpdatedAt        time.Time       `bun:"updated_at,notnull"`
}
