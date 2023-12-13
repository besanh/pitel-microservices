package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	ITemplateBss interface {
		IRepo[model.TemplateBss]
		GetTemplateBsses(ctx context.Context, db sqlclient.ISqlClientConn, filter model.TemplateBssFilter, limit, offset int) (total int, result *[]model.TemplateBssView, err error)
	}
	TemplateBss struct {
		Repo[model.TemplateBss]
	}
)

var TemplateBssRepo ITemplateBss

func NewTemplateBss() ITemplateBss {
	return &TemplateBss{}
}

func (repo *TemplateBss) GetTemplateBsses(ctx context.Context, db sqlclient.ISqlClientConn, filter model.TemplateBssFilter, limit, offset int) (total int, result *[]model.TemplateBssView, err error) {
	result = new([]model.TemplateBssView)
	query := db.GetDB().NewSelect().
		Model(result)
	if len(filter.TemplateName) > 0 {
		query.Where("? ILIKE %?%", bun.Ident("template_name"), filter.TemplateName)
	}

	if len(filter.TemplateCode) > 0 {
		query.Where("template_code IN (?)", bun.In(filter.TemplateCode))
	}

	if len(filter.TemplateType) > 0 {
		query.Where("template_type IN (?)", bun.In(filter.TemplateType))
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
