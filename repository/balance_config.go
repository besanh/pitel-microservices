package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IBalanceConfig interface {
		IRepo[model.BalanceConfig]
		GetBalanceConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.BalanceConfigFilter, limit, offset int) (total int, result *[]model.BalanceConfigView, err error)
		BulkInsertBalanceConfig(ctx context.Context, db sqlclient.ISqlClientConn, data []model.BalanceConfig) error
	}
	BalanceConfig struct {
		Repo[model.BalanceConfig]
	}
)

var BalanceConfigRepo IBalanceConfig

func NewBalanceConfig() IBalanceConfig {
	return &BalanceConfig{}
}

func (repo *BalanceConfig) BulkInsertBalanceConfig(ctx context.Context, db sqlclient.ISqlClientConn, data []model.BalanceConfig) error {
	query := db.GetDB().NewInsert().Model(&data)
	_, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BalanceConfig) GetBalanceConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.BalanceConfigFilter, limit, offset int) (total int, result *[]model.BalanceConfigView, err error) {
	result = new([]model.BalanceConfigView)
	query := db.GetDB().NewSelect().
		Model(result)
	if len(filter.Weight) > 0 {
		query.Where("weight IN (?)", bun.In(filter.Weight))
	}
	if len(filter.BalanceType) > 0 {
		query.Where("balance_type IN (?)", bun.In(filter.BalanceType))
	}
	if len(filter.Priority) > 0 {
		query.Where("priority IN (?)", bun.In(filter.Priority))
	}
	if len(filter.Provider) > 0 {
		query.Where("provider IN (?)", bun.In(filter.Provider))
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
