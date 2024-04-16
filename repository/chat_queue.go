package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatQueue interface {
		IRepo[model.ChatQueue]
		GetQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error)
	}
	ChatQueue struct {
		Repo[model.ChatQueue]
	}
)

var ChatQueueRepo IChatQueue

func NewChatQueue() IChatQueue {
	return &ChatQueue{}
}

func (repo *ChatQueue) GetById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatQueue, error) {
	result := new(model.ChatQueue)
	query := db.GetDB().NewSelect().Model(result).
		Relation("ConnectionQueues.ChatConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at desc")
		}).
		Relation("ChatRouting").
		Relation("ChatQueueUser").
		Relation("ChatManageQueueUser").
		Where("cq.id = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ChatQueue) GetQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error) {
	result := new([]model.ChatQueue)
	query := db.GetDB().NewSelect().Model(result).
		Relation("ConnectionQueues.ChatConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at desc")
		}).
		Relation("ChatRouting").
		Relation("ChatQueueUser").
		Relation("ChatManageQueueUser")
	if len(filter.TenantId) > 0 {
		query.Where("cq.tenant_id = ?", filter.TenantId)
	}
	if len(filter.QueueName) > 0 {
		query.Where("queue_name = ?", filter.QueueName)
	}
	if len(filter.QueueId) > 0 {
		query.Where("id = IN (?)", bun.In(filter.QueueId))
	}
	if len(filter.ChatRoutingId) > 0 {
		query.Where("chat_routing_id = IN (?)", bun.In(filter.ChatRoutingId))
	}
	query.Order("created_at desc")
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
