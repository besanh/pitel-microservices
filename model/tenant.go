package model

import "github.com/uptrace/bun"

type (
	BSS_Tenant struct {
		*Base
		bun.BaseModel `bun:"table:bss_tenants,alias:bt"`
		TenantName    string `bun:"tenant_name,type:text,notnull" json:"tenant_name"`
		Logo          string `bun:"logo,type:text" json:"logo"`
		Title         string `bun:"title,type:text" json:"title"`
		CssLink       string `bun:"css_link,type:text" json:"css_link"`
	}
	BSS_TenantRequest struct {
	}
)
