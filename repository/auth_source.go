package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IAuthSource interface {
		IRepo[model.AuthSource]
		GetAuthSource(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AuthSourceFilter, limit, offset int) (int, *[]model.AuthSource, error)
	}
	AuthSource struct {
		Repo[model.AuthSource]
	}
)

var AuthSourceRepo IAuthSource

func NewAuthSource() IAuthSource {
	return &AuthSource{}
}

func (s *AuthSource) GetAuthSource(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AuthSourceFilter, limit, offset int) (int, *[]model.AuthSource, error) {
	result := new([]model.AuthSource)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.Source) > 0 {
		query.Where("source = ?", filter.Source)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
