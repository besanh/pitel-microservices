package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IConnectionQueue interface {
		IRepo[model.ConnectionQueue]
		GetConnectionQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConnectionQueueFilter, limit, offset int) (int, *[]model.ConnectionQueue, error)
	}
	ConnectionQueue struct {
		Repo[model.ConnectionQueue]
	}
)

var ConnectionQueueRepo IConnectionQueue

func NewConnectionQueue() IConnectionQueue {
	return &ConnectionQueue{}
}

func (repo *ConnectionQueue) GetConnectionQueues(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConnectionQueueFilter, limit, offset int) (int, *[]model.ConnectionQueue, error) {
	result := new([]model.ConnectionQueue)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.ConnectionId) > 0 {
		query.Where("connection_id = ?", filter.ConnectionId)
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
