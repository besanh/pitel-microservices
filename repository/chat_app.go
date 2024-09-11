package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatApp interface {
		IRepo[model.ChatApp]
		GetChatApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAppFilter, limit, offset int) (int, *[]model.ChatApp, error)
	}
	ChatApp struct {
		Repo[model.ChatApp]
	}
)

var ChatAppRepo IChatApp

func NewChatApp() IChatApp {
	return &ChatApp{}
}

func (s *ChatApp) GetChatApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAppFilter, limit, offset int) (int, *[]model.ChatApp, error) {
	result := new([]model.ChatApp)
	query := db.GetDB().NewSelect().Model(result).
		Relation("ChatAppIntegrateSystems")
	if len(filter.AppName) > 0 {
		query.Where("app_name = ?", filter.AppName)
	}
	if len(filter.Status) > 0 {
		query.Where("status = ?", filter.Status)
	}
	if len(filter.AppType) > 0 {
		if len(filter.Status) > 0 {
			query.Where("info_app :: jsonb -> ? ->> 'status' = 'active'", filter.AppType)
		} else {
			query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("info_app :: jsonb -> ? ->> 'status' = 'active'", filter.AppType).
					WhereOr("info_app :: jsonb -> ? ->> 'status' = 'deactive'", filter.AppType)
			})
		}
	}
	if len(filter.AppId) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("info_app :: jsonb -> 'facebook' ->> 'app_id' = ?", filter.AppId).
				WhereOr("info_app :: jsonb -> 'zalo' ->> 'app_id' = ?", filter.AppId)
		})
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at desc")
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}

	return total, result, nil
}
