package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatManageQueue interface {
		IRepo[model.ChatManageQueueAgent]
		GetManageQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatManageQueueAgentFilter, limit, offset int) (int, *[]model.ChatManageQueueAgent, error)
	}
	ChatManageQueue struct {
		Repo[model.ChatManageQueueAgent]
	}
)

var ManageQueueRepo IChatManageQueue

func NewManageQueue() IChatManageQueue {
	return &ChatManageQueue{}
}

func (repo *ChatManageQueue) GetManageQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatManageQueueAgentFilter, limit, offset int) (int, *[]model.ChatManageQueueAgent, error) {
	entries := new([]model.ChatManageQueueAgent)
	query := db.GetDB().NewSelect().
		Model(entries).
		Limit(limit).
		Offset(offset)
	if len(filter.ConnectionId) > 0 {
		query.Where("connection_id = ?", filter.ConnectionId)
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.ManageId) > 0 {
		query.Where("manage_id = ?", filter.ManageId)
	}
	total, err := query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}
	return total, entries, nil
}
