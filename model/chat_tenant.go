package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatTenant struct {
	*Base
	bun.BaseModel     `bun:"table:chat_tenant,alias:ct"`
	TenantName        string               `json:"tenant_name" bun:"tenant_name,type:text,notnull"`
	IntegrateSystemId string               `json:"integrate_system_id" bun:"integrate_system_id,type:uuid,nullzero,default:null"`
	TenantId          string               `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	IntegrateSystem   *ChatIntegrateSystem `json:"integrate_system" bun:"rel:has-one,join:integrate_system_id=id"`
	Status            bool                 `json:"status" bun:"status,type:boolean,notnull"`
}

type ChatTenantRequest struct {
	IntegrateSystemId string `json:"integrate_system_id"`
	TenantName        string `json:"tenant_name" bind:"required"`
	Status            bool   `json:"status" bind:"required"`
}

func (m *ChatTenantRequest) Validate() error {
	if len(m.TenantName) < 1 {
		return errors.New("tenant name is required")
	}
	return nil
}
