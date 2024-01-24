package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
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

func (repo *ChatQueue) GetQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error) {
	result := new([]model.ChatQueue)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.QueueName) > 0 {
		query.Where("queue_name = ?", filter.QueueName)
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
