package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IRecipientConfig interface {
		IRepo[model.RecipientConfig]
		BulkInsertRecipientConfig(ctx context.Context, db sqlclient.ISqlClientConn, data []model.RecipientConfig) error
		GetRecipientConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.RecipientConfigFilter, limit, offset int) (total int, result *[]model.RecipientConfigView, err error)
	}
	RecipientConfig struct {
		Repo[model.RecipientConfig]
	}
)

var RecipientConfigRepo IRecipientConfig

func NewRecipientConfig() IRecipientConfig {
	return &RecipientConfig{}
}

func (repo *RecipientConfig) BulkInsertRecipientConfig(ctx context.Context, db sqlclient.ISqlClientConn, data []model.RecipientConfig) error {
	query := db.GetDB().NewInsert().Model(&data)
	_, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *RecipientConfig) GetRecipientConfigs(ctx context.Context, db sqlclient.ISqlClientConn, filter model.RecipientConfigFilter, limit, offset int) (total int, result *[]model.RecipientConfigView, err error) {
	result = new([]model.RecipientConfigView)
	query := db.GetDB().NewSelect().
		Model(result)
	if len(filter.Recipient) > 0 {
		query.Where("recipient IN (?)", bun.In(filter.Recipient))
	}
	if len(filter.RecipientType) > 0 {
		query.Where("recipient_type IN (?)", bun.In(filter.RecipientType))
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
