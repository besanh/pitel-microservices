package model

import (
	"github.com/uptrace/bun"
	"time"
)

type ChatAutoScriptToChatScript struct {
	bun.BaseModel    `bun:"table:chat_auto_script_to_chat_script,alias:cas_cst"`
	ChatAutoScriptId string          `json:"chat_auto_script_id" bun:"chat_auto_script_id,type:uuid,pk"`
	ChatScriptId     string          `json:"chat_script_id" bun:"chat_script_id,type:uuid,pk"`
	Order            int             `json:"order" bun:"order,notnull"`
	ChatAutoScript   *ChatAutoScript `json:"chat_auto_script" bun:"rel:belongs-to,join:chat_auto_script_id=id"`
	ChatScript       *ChatScript     `json:"chat_script" bun:"rel:belongs-to,join:chat_script_id=id"`
	CreatedAt        time.Time       `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt        time.Time       `json:"updated_at" bun:"updated_at,notnull"`
}
