package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatPolicySetting interface {
		IRepo[model.ChatPolicySetting]
		GetChatPolicySettings(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatPolicyFilter, limit, offset int) (int, *[]model.ChatPolicySetting, error)
	}

	ChatPolicySetting struct {
		Repo[model.ChatPolicySetting]
	}
)

var ChatPolicySettingRepo IChatPolicySetting

func NewChatPolicySetting() IChatPolicySetting {
	return &ChatPolicySetting{}
}

func (c ChatPolicySetting) GetChatPolicySettings(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatPolicyFilter, limit, offset int) (int, *[]model.ChatPolicySetting, error) {
	result := new([]model.ChatPolicySetting)
	query := db.GetDB().NewSelect().Model(result).
		Column("cps.*")
	if len(filter.TenantId) > 0 {
		query.Where("cps.tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConnectionType) > 0 {
		query.Where("cps.connection_type = ?", filter.ConnectionType)
	}
	if len(filter.ExcludedIds) > 0 {
		query.Where("cps.id NOT IN (?)", bun.In(filter.ExcludedIds))
	}

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("cps.created_at desc")

	total, err := query.ScanAndCount(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, result, nil
	} else if err != nil {
		return 0, nil, err
	}
	return total, result, nil
}
