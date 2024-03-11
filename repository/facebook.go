package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IFacebook interface {
		IRepo[model.FacebookPage]
		GetFacebookPages(ctx context.Context, db sqlclient.ISqlClientConn, filter model.FacebookPageFilter, limit, offset int) (int, *[]model.FacebookPage, error)
	}
	Facebook struct {
		Repo[model.FacebookPage]
	}
)

func NewFacebook() IFacebook {
	repo := &Facebook{}
	return repo
}

func (repo *Facebook) GetFacebookPages(ctx context.Context, db sqlclient.ISqlClientConn, filter model.FacebookPageFilter, limit, offset int) (int, *[]model.FacebookPage, error) {
	result := new([]model.FacebookPage)
	query := db.GetDB().NewSelect().
		Model(result)
	if len(filter.OaId) > 0 {
		query.Where("oa_id = ?", filter.OaId)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at desc")

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
