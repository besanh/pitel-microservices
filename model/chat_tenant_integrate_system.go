package model

import "github.com/uptrace/bun"

type ChatTenantIntegrateSystem struct {
	*Base
	bun.BaseModel       `bun:"table:chat_tenant_integrate_system,alias:ctis"`
	TenantId            string               `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	IntegrateSystemId   string               `json:"integrate_system_id" bun:"integrate_system_id,type:uuid,notnull"`
	ChatIntegrateSystem *ChatIntegrateSystem `json:"chat_integrate_system" bun:"chat_integrate_system,type:jsonb,notnull"`
	Status              bool                 `json:"status" bun:"status,type:boolean,notnull"`
}
