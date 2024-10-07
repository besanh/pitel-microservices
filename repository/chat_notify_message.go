package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatNotifyMessage interface {
		IRepo[model.ChatNotifyMessage]
		GetChatNotifyMessages(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatNotifyMessageFilter, limit, offset int) (total int, entries *[]model.ChatNotifyMessage, err error)
		GetChatNotifyMessageById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (result *model.ChatNotifyMessage, err error)
	}
	ChatNotifyMessage struct {
		Repo[model.ChatNotifyMessage]
	}
)

func (repo *ChatNotifyMessage) GetChatNotifyMessageById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (result *model.ChatNotifyMessage, err error) {
	result = &model.ChatNotifyMessage{}
	err = db.GetDB().NewSelect().
		Model(result).
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name", "connection_type", "oa_info")
		}).
		Where("cnm.id = ?", id).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return
	}
	return
}

var ChatNotifyMessageRepo IChatNotifyMessage

func NewChatNotifyMessage() IChatNotifyMessage {
	return &ChatNotifyMessage{}
}

func (repo *ChatNotifyMessage) GetChatNotifyMessages(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatNotifyMessageFilter, limit, offset int) (total int, entries *[]model.ChatNotifyMessage, err error) {
	entries = new([]model.ChatNotifyMessage)
	query := db.GetDB().NewSelect().
		Model(entries).
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name", "connection_type", "oa_info")
		})
	if len(filter.TenantId) > 0 {
		query.Where("cnm.tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConnectionId) > 0 {
		query.Where("cnm.connection_id = ?", filter.ConnectionId)
	}
	if len(filter.NotifyType) > 0 {
		query.Where("cnm.notify_type IN (?)", bun.In(filter.NotifyType))
	}
	query.Order("cnm.created_at desc")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, nil, err
	}

	return total, entries, nil
}
