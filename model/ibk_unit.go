package model

import (
	"github.com/uptrace/bun"
)

type (
	IBKBusinessUnit struct {
		*Base
		bun.BaseModel    `bun:"table:ibk_business_units"`
		TenantId         string `bun:"tenant_id,type:uuid,notnull"`
		ParentId         string `bun:"parent_id,type:uuid,nullzero,default:null"`
		BusinessUnitName string `bun:"business_unit_name,type:text,notnull"`
		Address          string `bun:"address,type:text"`
	}
	IBKBusinessUnitInfo struct {
		*IBKBusinessUnit
		bun.BaseModel `bun:"table:ibk_business_units,alias:bu"`
		TenantName    string `bun:"tenant_name"`
		TotalUser     int32  `bun:"total_user"`
	}
)
type (
	IBKBusinessUnitQueryParam struct {
		TenantId string
		ParentId string
	}
	IBKBusinessUnitBody struct {
		BusinessUnitName string `json:"business_unit_name" required:"true" pattern:"^[a-zA-Z0-9 _-]{0,50}$" doc:"Bussiness Unit Name"`
		ParentId         string `json:"parent_id" required:"true" pattern:"^[a-zA-Z0-9 _-]{0,50}$" doc:"Parent Id"`
		Address          string `json:"address" required:"false"  nullable:"true" format:"uuid" doc:"Address"`
	}
)
