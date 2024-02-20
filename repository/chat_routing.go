package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatRouting interface {
		IRepo[model.ChatRouting]
		GetChatRoutings(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatRoutingFilter, limit, offset int) (int, *[]model.ChatRouting, error)
	}
	ChatRouting struct {
		Repo[model.ChatRouting]
	}
)

var ChatRoutingRepo IChatRouting

func NewChatRouting() IChatRouting {
	return &ChatRouting{}
}

func (repo *ChatRouting) GetChatRoutings(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatRoutingFilter, limit, offset int) (int, *[]model.ChatRouting, error) {
	result := new([]model.ChatRouting)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.RoutingName) > 0 {
		query.Where("routing_name = ?", filter.RoutingName)
	}
	if len(filter.RoutingAlias) > 0 {
		query.Where("routing_alias = ?", filter.RoutingAlias)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
