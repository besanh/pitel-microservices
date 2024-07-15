package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ChatAppIntegrateSystem struct {
	bun.BaseModel         `bun:"table:chat_app_integrate_system,alias:cais"`
	ChatAppId             string               `json:"chat_app_id" bun:"chat_app_id,pk,type:uuid"`
	ChatApp               *ChatApp             `json:"chat_app" bun:"rel:belongs-to,join:chat_app_id=id"`
	ChatIntegrateSystemId string               `json:"chat_integrate_system_id" bun:"chat_integrate_system_id,pk,type:uuid"`
	ChatIntegrateSystem   *ChatIntegrateSystem `json:"chat_integrate_system" bun:"rel:belongs-to,join:chat_integrate_system_id=id"`
	CreatedAt             time.Time            `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt             time.Time            `json:"updated_at" bun:"updated_at"`
}
