package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IShareInfo interface {
		IRepo[model.ShareInfoForm]
		GetShareInfos(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ShareInfoFormFilter, limit, offset int) (int, *[]model.ShareInfoForm, error)
	}
	ShareInfo struct {
		Repo[model.ShareInfoForm]
	}
)

var ShareInfoRepo IShareInfo

func NewShareInfo() IShareInfo {
	return &ShareInfo{}
}

func (s *ShareInfo) GetShareInfos(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ShareInfoFormFilter, limit, offset int) (int, *[]model.ShareInfoForm, error) {
	result := new([]model.ShareInfoForm)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.OaId) > 0 {
		query.Where("share_form->?::text->'oa_id' = ?", filter.ShareType, filter.OaId)
	}
	if len(filter.AppId) > 0 {
		query.Where("share_form->?->>'app_id' = ?", filter.ShareType, filter.AppId)
	}
	if len(filter.UserId) > 0 {
		query.Where("user_id = ?", filter.UserId)
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
