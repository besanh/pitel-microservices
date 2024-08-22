package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatTenant interface {
		IRepo[model.ChatTenant]
		GetChatTenants(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatTenantFilter, limit, offset int) (total int, result *[]model.ChatTenant, err error)
		GetByTenantId(ctx context.Context, db sqlclient.ISqlClientConn, tenantId string) (entity *model.ChatTenant, err error)
		DeleteByTenantId(ctx context.Context, db sqlclient.ISqlClientConn, tenantId string) (err error)
	}
	ChatTenant struct {
		Repo[model.ChatTenant]
	}
)

var ChatTenantRepo IChatTenant

func NewChatTenant() IChatTenant {
	return &ChatTenant{}
}

func (repo *ChatTenant) GetChatTenants(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatTenantFilter, limit, offset int) (total int, result *[]model.ChatTenant, err error) {
	result = new([]model.ChatTenant)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantName) > 0 {
		query.Where("tenant_name = ?", filter.TenantName)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}

func (repo *ChatTenant) GetByTenantId(ctx context.Context, db sqlclient.ISqlClientConn, tenantId string) (entity *model.ChatTenant, err error) {
	entity = new(model.ChatTenant)
	err = db.GetDB().NewSelect().
		Model(entity).
		Where("tenant_id = ?", tenantId).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ChatTenant) DeleteByTenantId(ctx context.Context, db sqlclient.ISqlClientConn, tenantId string) (err error) {
	_, err = db.GetDB().NewDelete().
		Model((*model.ChatTenant)(nil)).
		Where("tenant_id = ?", tenantId).
		Exec(ctx)
	return
}
