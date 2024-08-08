package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IIBKTenant interface {
		IRepo[model.IBKTenant]
		GetInfoByQuery(c context.Context, param model.IBKTenantQueryParam, limit, offset int) (total int, result []*model.IBKTenantInfo, err error)
		GetInfoById(c context.Context, id string) (result *model.IBKTenantInfo, err error)
	}
	IBKTenant struct {
		Repo[model.IBKTenant]
	}
)

var IBKTenantRepo IIBKTenant

func NewIBKTenant(conn sqlclient.ISqlClientConn) IIBKTenant {
	return &IBKTenant{
		Repo[model.IBKTenant]{
			Conn: conn,
		},
	}
}

func (repo *IBKTenant) GetInfoByQuery(c context.Context, param model.IBKTenantQueryParam, limit, offset int) (total int, result []*model.IBKTenantInfo, err error) {
	result = make([]*model.IBKTenantInfo, 0)
	query := repo.Conn.GetDB().NewSelect().
		Model(&result).
		Limit(limit).
		Offset(offset)
	if len(param.TenantId_Eq) > 0 {
		query.Where("tenant_id = ?", param.TenantId_Eq)
	}
	query.ColumnExpr("tenant.*").
		ColumnExpr("(?) as total_business_unit", repo.Conn.GetDB().NewSelect().TableExpr("ibk_business_units as bu").
			ColumnExpr("count(bu.id) as total_business_unit").
			Where("bu.tenant_id = tenant.id"),
		).
		ColumnExpr("(?) as total_user", repo.Conn.GetDB().NewSelect().TableExpr("ibk_users as u").
			ColumnExpr("count(u.id) as total_user").
			Join("INNER JOIN ibk_business_unit bu ON g.business_unit_id = bu.id").
			Where("bu.tenant_id = tenant.id"),
		)
	total, err = query.ScanAndCount(c)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}

func (repo *IBKTenant) GetInfoById(c context.Context, id string) (result *model.IBKTenantInfo, err error) {
	result = new(model.IBKTenantInfo)
	err = repo.Conn.GetDB().NewSelect().Model(result).
		ColumnExpr("tenant.*").
		ColumnExpr("(?) as total_business_unit", repo.Conn.GetDB().NewSelect().TableExpr("ibk_business_units as bu").
			ColumnExpr("count(bu.id) as total_business_unit").
			Where("bu.tenant_id = tenant.id"),
		).
		ColumnExpr("(?) as total_user", repo.Conn.GetDB().NewSelect().TableExpr("ibk_users as u").
			ColumnExpr("count(u.id) as total_user").
			Join("INNER JOIN ibk_business_unit bu ON g.business_unit_id = bu.id").
			Where("bu.tenant_id = tenant.id"),
		).
		Where("tenant.id = ?", id).
		Scan(c)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}
