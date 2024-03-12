package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IAgentAllocation interface {
		IRepo[model.AgentAllocation]
		GetAgentAllocations(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AgentAllocationFilter, limit, offset int) (int, *[]model.AgentAllocation, error)
	}
	AgentAllocation struct {
		Repo[model.AgentAllocation]
	}
)

var AgentAllocationRepo IAgentAllocation

func NewAgentAllocation() IAgentAllocation {
	return &AgentAllocation{}
}

func (repo *AgentAllocation) GetAgentAllocations(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AgentAllocationFilter, limit, offset int) (int, *[]model.AgentAllocation, error) {
	result := new([]model.AgentAllocation)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.AppId) > 0 {
		query.Where("app_id = ?", filter.AppId)
	}
	if len(filter.AgentId) > 0 {
		query.Where("agent_id IN (?)", bun.In(filter.AgentId))
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.ConversationId) > 0 {
		query.Where("conversation_id = ?", filter.ConversationId)
	}
	if len(filter.MainAllocate) > 0 {
		query.Where("main_allocate = ?", filter.MainAllocate)
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
