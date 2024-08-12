package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IIBKUser interface {
		IRepo[model.IBKUser]
		GetUserByUsername(ctx context.Context, userName string) (result *model.IBKUser, err error)
		GetInfoById(ctx context.Context, id string, params model.IBKUserQueryParam) (user *model.IBKUserInfo, err error)
		GetInfoByQuery(ctx context.Context, params model.IBKUserQueryParam, limit, offset int, orderBy ...string) (total int, result []*model.IBKUserInfo, err error)
	}
	IBKUser struct {
		Repo[model.IBKUser]
	}
)

var IBKUserRepo IIBKUser

func NewIBKUser(conn sqlclient.ISqlClientConn) IIBKUser {
	return &IBKUser{
		Repo[model.IBKUser]{
			Conn: conn,
		},
	}
}

func (repo *IBKUser) GetUserByUsername(ctx context.Context, userName string) (result *model.IBKUser, err error) {
	result = new(model.IBKUser)
	query := repo.Conn.GetDB().NewSelect().
		ColumnExpr("*").
		Model(result).
		Where("username = ?", userName).
		Limit(1)
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		err = nil
		result = nil
	}
	return
}

func (repo *IBKUser) GetInfoByQuery(ctx context.Context, params model.IBKUserQueryParam, limit, offset int, orderBy ...string) (total int, result []*model.IBKUserInfo, err error) {
	result = make([]*model.IBKUserInfo, 0)
	query := repo.Conn.GetDB().NewSelect().Model(&result).
		ColumnExpr("DISTINCT u.id").
		ColumnExpr("u.*").
		ColumnExpr("bu.id as business_unit_id, bu.business_unit_name").
		ColumnExpr("t.id as tenant_id, t.tenant_name").
		ColumnExpr("r.id as role_id, r.role_name").
		Join("JOIN ibk_business_units bu ON u.business_unit_id = bu.id").
		Join("JOIN ibk_tenants t ON bu.tenant_id = t.id").
		Join("JOIN ibk_roles r ON r.id = u.role_id").
		Offset(offset).
		Limit(limit)
	if len(params.BusinessUnitId_Eq) > 0 {
		query.Where("u.business_unit_id::text = ?", params.BusinessUnitId_Eq)
	}
	if len(params.TenantId_Eq) > 0 {
		query.Where("t.id::text = ?", params.TenantId_Eq)
	}
	if len(params.Keyword) > 0 {
		query.Where("u.username LIKE ? OR u.fullname LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if len(params.RoleId_Eq) > 0 || len(params.ServiceId_Eq) > 0 {
		subQuery := repo.Conn.GetDB().NewSelect().
			TableExpr("ibk_user_service as us").
			Where("us.user_id = u.id")
		if len(params.RoleId_Eq) > 0 {
			subQuery.Where("us.role_id::text = ?", params.RoleId_Eq)
		}
		if len(params.ServiceId_Eq) > 0 {
			subQuery.Where("us.service_id::text = ?", params.ServiceId_Eq)
		}
		query.Where("EXISTS (?)", subQuery)
	}
	if len(params.Username_Eq) > 0 {
		query.Where("u.username = ?", params.Username_Eq)
	}
	if len(params.Fullname_Eq) > 0 {
		query.Where("u.fullname = ?", params.Fullname_Eq)
	}
	if len(params.Fullname_Like) > 0 {
		query.Where("u.fullname LIKE ?", "%"+params.Fullname_Like+"%")
	}
	if len(params.Email_Eq) > 0 {
		query.Where("u.email = ?", params.Email_Eq)
	}
	if len(params.Email_Like) > 0 {
		query.Where("u.email LIKE ?", "%"+params.Email_Like+"%")
	}
	// if params.IsActivated_Eq != nil {
	// 	query.Where("u.is_activated = ?", params.IsActivated_Eq)
	// }
	// if params.IsLocked_Eq != nil {
	// 	query.Where("u.is_locked = ?", params.IsLocked_Eq)
	// }
	orderValue := "u.created_at"
	sortValue := "DESC"
	switch params.Order {
	case "created_at":
		orderValue = "u.created_at"
	case "updated_at":
		orderValue = "u.updated_at"
	case "fullname":
		orderValue = "u.fullname"
	case "username":
		orderValue = "u.username"
	case "business_unit_id":
		orderValue = "bu.id"
	}
	if len(params.Sort) > 0 && params.Sort == "asc" {
		sortValue = "ASC"
	}
	query.OrderExpr(orderValue + " " + sortValue)

	if total, err = query.ScanAndCount(ctx); err == sql.ErrNoRows {
		result = nil
		return
	}
	return
}

func (repo *IBKUser) GetInfoById(ctx context.Context, id string, params model.IBKUserQueryParam) (user *model.IBKUserInfo, err error) {
	user = new(model.IBKUserInfo)
	query := repo.Conn.GetDB().NewSelect().Model(user).
		ColumnExpr("u.*").
		ColumnExpr("bu.id as business_unit_id, bu.business_unit_name").
		ColumnExpr("t.id as tenant_id, t.tenant_name").
		ColumnExpr("r.id as role_id, r.role_name").
		Join("JOIN ibk_business_units bu ON u.business_unit_id = bu.id").
		Join("JOIN ibk_tenants t ON bu.tenant_id = t.id").
		Join("JOIN ibk_roles r ON r.id = u.role_id").
		Where("u.id::text = ?", id).
		Limit(1)

	if len(params.BusinessUnitId_Eq) > 0 {
		query.Where("u.business_unit_id::text = ?", params.BusinessUnitId_Eq)
	}
	if len(params.TenantId_Eq) > 0 {
		query.Where("t.id::text = ?", params.TenantId_Eq)
	}
	if len(params.Keyword) > 0 {
		query.Where("u.username LIKE ? OR u.fullname LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if len(params.RoleId_Eq) > 0 || len(params.ServiceId_Eq) > 0 {
		subQuery := repo.Conn.GetDB().NewSelect().
			TableExpr("IBKUser_service as us").
			Where("us.user_id = u.id")
		if len(params.RoleId_Eq) > 0 {
			subQuery.Where("us.role_id::text = ?", params.RoleId_Eq)
		}
		if len(params.ServiceId_Eq) > 0 {
			subQuery.Where("us.service_id::text = ?", params.ServiceId_Eq)
		}
		query.Where("EXISTS (?)", subQuery)
	}
	if len(params.Username_Eq) > 0 {
		query.Where("u.username = ?", params.Username_Eq)
	}
	if len(params.Email_Eq) > 0 {
		query.Where("u.email = ?", params.Email_Eq)
	}
	if err = query.Scan(ctx); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return
}
