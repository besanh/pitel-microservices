package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IPluginConfig interface {
		IRepo[model.PluginConfig]
		GetPluginConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.PluginConfigFilter, limit, offset int) (int, *[]model.PluginConfigView, error)
	}
	PluginConfig struct {
		Repo[model.PluginConfig]
	}
)

var PluginConfigRepo IPluginConfig

func NewPluginConfig() IPluginConfig {
	return &PluginConfig{}
}

func (repo *PluginConfig) GetPluginConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.PluginConfigFilter, limit, offset int) (int, *[]model.PluginConfigView, error) {
	result := new([]model.PluginConfigView)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.PluginName) > 0 {
		query.Where("plugin_name IN (?)", bun.In(filter.PluginName))
	}
	if len(filter.PluginType) > 0 {
		query.Where("plugin_type IN (?)", bun.In(filter.PluginType))
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status)
	}
	query.Order("created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return total, result, nil
	} else if err != nil {
		return 0, result, err
	}

	return total, result, nil
}
