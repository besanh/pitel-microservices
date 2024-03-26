package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatQueueUser interface {
		IRepo[model.ChatQueueUser]
		GetChatQueueUsers(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatQueueUserFilter, limit, offset int) (int, *[]model.ChatQueueUser, error)
		DeleteChatQueueUsers(ctx context.Context, db sqlclient.ISqlClientConn, queueId string) error
	}
	ChatQueueUser struct {
		Repo[model.ChatQueueUser]
	}
)

var ChatQueueUserRepo IChatQueueUser

func NewChatQueueUser() IChatQueueUser {
	return &ChatQueueUser{}
}

func (repo *ChatQueueUser) GetChatQueueUsers(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatQueueUserFilter, limit, offset int) (int, *[]model.ChatQueueUser, error) {
	result := new([]model.ChatQueueUser)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id IN (?)", bun.In(filter.QueueId))
	}
	if len(filter.UserId) > 0 {
		query.Where("user_id IN (?)", bun.In(filter.UserId))
	}
	if len(filter.Source) > 0 {
		query.Where("source = ?", filter.Source)
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

func (repo *ChatQueueUser) DeleteChatQueueUsers(ctx context.Context, db sqlclient.ISqlClientConn, queueId string) error {
	_, err := db.GetDB().NewDelete().Model(new(model.ChatQueueUser)).Where("queue_id = ?", queueId).Exec(ctx)
	return err
}
