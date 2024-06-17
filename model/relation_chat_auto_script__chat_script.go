package model

import (
	"github.com/uptrace/bun"
	"time"
)

type ChatAutoScriptToChatScript struct {
	bun.BaseModel    `bun:"table:chat_auto_script_to_chat_script,alias:caslink"`
	ChatAutoScriptId string          `bun:"chat_auto_script_id,type:uuid,pk"`
	ChatScriptId     string          `bun:"chat_script_id,type:uuid,pk"`
	Order            int             `bun:"order,notnull"`
	ChatAutoScript   *ChatAutoScript `bun:"rel:belongs-to,join:chat_auto_script_id=id"`
	ChatScript       *ChatScript     `bun:"rel:belongs-to,join:chat_script_id=id"`
	CreatedAt        time.Time       `bun:"created_at,notnull"`
	UpdatedAt        time.Time       `bun:"updated_at,notnull"`
}
