package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IAgentAllocation interface {
		IRepo[model.AgentAllocation]
		GetAgentAllocation(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AgentAllocationFilter, limit, offset int) (int, *[]model.AgentAllocation, error)
	}
	AgentAllocation struct {
		Repo[model.AgentAllocation]
	}
)

var AgentAllocationRepo IAgentAllocation

func NewAgentAllocation() IAgentAllocation {
	return &AgentAllocation{}
}

func (repo *AgentAllocation) GetAgentAllocation(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AgentAllocationFilter, limit, offset int) (int, *[]model.AgentAllocation, error) {
	result := new([]model.AgentAllocation)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.UserIdByApp) > 0 {
		query.Where("user_id_by_app = ?", filter.UserIdByApp)
	}
	if len(filter.AgentId) > 0 {
		query.Where("agent_id = ?", filter.AgentId)
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
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
