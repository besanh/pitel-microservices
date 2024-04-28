package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatManageQueue interface {
		IRepo[model.ChatManageQueueUser]
		GetManageQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatManageQueueUserFilter, limit, offset int) (int, *[]model.ChatManageQueueUser, error)
	}
	ChatManageQueue struct {
		Repo[model.ChatManageQueueUser]
	}
)

var ManageQueueRepo IChatManageQueue

func NewManageQueue() IChatManageQueue {
	return &ChatManageQueue{}
}

func (repo *ChatManageQueue) GetManageQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatManageQueueUserFilter, limit, offset int) (int, *[]model.ChatManageQueueUser, error) {
	entries := new([]model.ChatManageQueueUser)
	query := db.GetDB().NewSelect().
		Model(entries)
	if len(filter.ConnectionId) > 0 {
		query.Where("connection_id = ?", filter.ConnectionId)
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.ManageId) > 0 {
		query.Where("manage_id = ?", filter.ManageId)
	}
	query.Order("created_at desc")
	if limit > 0 {
		query.Limit(limit).
			Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}
	return total, entries, nil
}
