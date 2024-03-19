package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IManageQueue interface {
		IRepo[model.ChatManageQueueAgent]
		GetManageQueue(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ManageQueueAgentFilter, limit, offset int) (int, *[]model.ChatManageQueueAgent, error)
	}
	ManageQueue struct {
		Repo[model.ChatManageQueueAgent]
	}
)

var ManageQueueRepo IManageQueue

func NewManageQueue() IManageQueue {
	return &ManageQueue{}
}

func (repo *ManageQueue) GetManageQueue(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ManageQueueAgentFilter, limit, offset int) (int, *[]model.ChatManageQueueAgent, error) {
	entries := new([]model.ChatManageQueueAgent)
	query := db.GetDB().NewSelect().
		Model(entries).
		Limit(limit).
		Offset(offset)
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.AgentId) > 0 {
		query.Where("agent_id = ?", filter.AgentId)
	}
	total, err := query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}
	return total, entries, nil
}
