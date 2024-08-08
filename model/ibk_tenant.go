package model

import (
	"github.com/uptrace/bun"
)

type (
	IBKTenant struct {
		*Base
		bun.BaseModel `bun:"table:ibk_tenants"`
		TenantName    string         `bun:"tenant_name,type:text,notnull" json:"tenant_name"`
		Logo          string         `bun:"logo,type:text" json:"logo"`
		MetaData      map[string]any `bun:"meta_data,type:jsonb" json:"meta_data"`
	}
	IBKTenantInfo struct {
		*IBKTenant
		bun.BaseModel     `bun:"table:ibk_tenant,alias:tenant"`
		TotalBusinessUnit int32 `bun:"total_business_unit,nullzero"`
		TotalUser         int32 `bun:"total_user,nullzero"`
	}
)
type (
	IBKTenantQueryParam struct {
		Keyword     string `query:"keyword"`
		Sort        string `query:"sort"`
		Order       string `query:"order"`
		TenantId_Eq string `query:"tenant_id"`
	}
	IBKTenantBody struct {
		TenantName string `json:"tenant_name" required:"true"`
	}
)
