package repository

import (
	"context"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
)

type (
	IChatAppIntegrateSystem interface {
		IRepo[model.ChatAppIntegrateSystem]
		GetChatAppIntegrateSystems(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAppIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatAppIntegrateSystem, err error)
	}
	ChatAppIntegrateSystem struct {
		Repo[model.ChatAppIntegrateSystem]
	}
)

var ChatAppIntegrateSystemRepo IChatAppIntegrateSystem

func NewChatAppIntegrateSystem() IChatAppIntegrateSystem {
	return &ChatAppIntegrateSystem{}
}

func (repo *ChatAppIntegrateSystem) GetChatAppIntegrateSystems(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAppIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatAppIntegrateSystem, err error) {
	result = new([]model.ChatAppIntegrateSystem)
	query := db.GetDB().NewSelect().Model(result).
		Relation("ChatIntegrateSystem")
	if len(filter.ChatAppId) > 0 {
		query.Where("chat_app_id = ?", filter.ChatAppId)
	}
	if len(filter.ChatIntegrateSystemId) > 0 {
		query.Where("chat_integrate_system_id = ?", filter.ChatIntegrateSystemId)
	}
	query.Order("created_at DESC")
	if limit > 0 {
		query.Limit(limit).
			Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err != nil {
		return 0, nil, err
	}
	return total, result, nil
}
