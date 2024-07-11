package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatIntegrateSystem interface {
		IRepo[model.ChatIntegrateSystem]
		GetIntegrateSystem(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error)
	}
	ChatIntegrateSystem struct {
		Repo[model.ChatIntegrateSystem]
	}
)

var ChatIntegrateSystemRepo IChatIntegrateSystem

func NewChatIntegrateSystem() IChatIntegrateSystem {
	return &ChatIntegrateSystem{}
}

func (repo *ChatIntegrateSystem) GetIntegrateSystem(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error) {
	result = new([]model.ChatIntegrateSystem)
	query := db.GetDB().NewSelect().Model(result).
		Relation("Vendor", func(q *bun.SelectQuery) *bun.SelectQuery {
			if len(filter.VendorName) > 0 {
				return q.Where("vendor_name = ?", filter.VendorName)
			}
			return q
		})
	if len(filter.SystemName) > 0 {
		query.Where("system_name = ?", filter.SystemName)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status)
	}
	query.Order("created_at DESC")

	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
