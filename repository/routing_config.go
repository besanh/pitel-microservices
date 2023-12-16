package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IRoutingConfig interface {
		IRepo[model.RoutingConfig]
		GetRoutingConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.RoutingConfigFilter, limit, offset int) (total int, result *[]model.RoutingConfigView, err error)
	}
	RoutingConfig struct {
		Repo[model.RoutingConfig]
	}
)

var RoutingConfigRepo IRoutingConfig

func NewRoutingConfig() IRoutingConfig {
	return &RoutingConfig{}
}

func (repo *RoutingConfig) GetRoutingConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.RoutingConfigFilter, limit, offset int) (total int, result *[]model.RoutingConfigView, err error) {
	result = new([]model.RoutingConfigView)
	query := db.GetDB().NewSelect().Model(&result)
	if len(filter.RoutingName) > 0 {
		query.Where("routing_name = ?", filter.RoutingName)
	}
	if len(filter.RoutingType) > 0 {
		query.Where("routing_type IN (?)", bun.In(filter.RoutingType))
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	query.Order("created_at DESC")

	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, nil
	} else if err != nil {
		return 0, nil, err
	}
	return total, result, nil
}
