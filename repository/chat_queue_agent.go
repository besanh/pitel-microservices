package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatQueueAgent interface {
		IRepo[model.ChatQueueAgent]
		GetChatQueueAgents(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatQueueAgentFilter, limit, offset int) (int, *[]model.ChatQueueAgent, error)
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
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.AgentId) > 0 {
		query.Where("agent_id = ?", filter.AgentId)
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
