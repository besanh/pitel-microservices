package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IIBKRole interface {
		IRepo[model.IBKRole]
		SelectInfo(ctx context.Context, params model.IBKRoleQueryParam, limit int, offset int) (entries []*model.IBKRoleInfo, total int, err error)
		FindInfoById(ctx context.Context, id string) (entry *model.IBKRoleInfo, err error)
	}
	IBKRole struct {
		Repo[model.IBKRole]
	}
)

var IBKRoleRepo IIBKRole

func NewIBKRole(conn sqlclient.ISqlClientConn) IIBKRole {
	repo := &IBKRole{
		Repo[model.IBKRole]{
			Conn: conn,
		},
	}
	return repo
}

func (repo *IBKRole) SelectInfo(ctx context.Context, params model.IBKRoleQueryParam, limit int, offset int) (entries []*model.IBKRoleInfo, total int, err error) {
	entries = make([]*model.IBKRoleInfo, 0)
	query := repo.Conn.GetDB().NewSelect().
		Model(&entries).
		Limit(limit).
		Offset(offset)

	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (repo *IBKRole) FindInfoById(ctx context.Context, id string) (entry *model.IBKRoleInfo, err error) {
	entry = new(model.IBKRoleInfo)
	err = repo.Conn.GetDB().NewSelect().
		Model(entry).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
