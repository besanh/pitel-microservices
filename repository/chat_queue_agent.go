package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatQueueAgent interface {
		IRepo[model.ChatQueueAgent]
		GetChatQueueAgents(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatQueueAgentFilter, limit, offset int) (int, *[]model.ChatQueueAgent, error)
		DeleteChatQueueAgents(ctx context.Context, db sqlclient.ISqlClientConn, queueId string) error
	}
	ChatQueueAgent struct {
		Repo[model.ChatQueueAgent]
	}
)

var ChatQueueAgentRepo IChatQueueAgent

func NewChatQueueAgent() IChatQueueAgent {
	return &ChatQueueAgent{}
}

func (repo *ChatQueueAgent) GetChatQueueAgents(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatQueueAgentFilter, limit, offset int) (int, *[]model.ChatQueueAgent, error) {
	result := new([]model.ChatQueueAgent)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id IN (?)", bun.In(filter.QueueId))
	}
	if len(filter.AgentId) > 0 {
		query.Where("agent_id IN (?)", bun.In(filter.AgentId))
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

func (repo *ChatQueueAgent) DeleteChatQueueAgents(ctx context.Context, db sqlclient.ISqlClientConn, queueId string) error {
	_, err := db.GetDB().NewDelete().Model(new(model.ChatQueueAgent)).Where("queue_id = ?", queueId).Exec(ctx)
	return err
}
