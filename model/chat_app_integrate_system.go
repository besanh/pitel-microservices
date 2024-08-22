package model

import (
	"github.com/uptrace/bun"
)

type ChatAppIntegrateSystem struct {
	*Base
	bun.BaseModel         `bun:"table:chat_app_integrate_system,alias:cais"`
	ChatAppId             string                 `json:"chat_app_id" bun:"chat_app_id,type:uuid,notnull"`
	ChatApp               *ChatApp               `json:"chat_app" bun:"rel:belongs-to,join:chat_app_id=id"`
	ChatIntegrateSystemId string                 `json:"chat_integrate_system_id" bun:"chat_integrate_system_id,type:uuid,notnull"`
	ChatIntegrateSystem   []*ChatIntegrateSystem `json:"chat_integrate_system" bun:"rel:has-many,join:chat_integrate_system_id=id"`
}
