package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IConnectionApp interface {
		IRepo[model.ConnectionApp]
		GetConnectionApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConnectionAppFilter, limit, offset int) (int, *[]model.ConnectionApp, error)
	}
	ConnectionApp struct {
		Repo[model.ConnectionApp]
	}
)

var ConnectionAppRepo IConnectionApp

func NewConnectionApp() IConnectionApp {
	return &ConnectionApp{}
}

func (repo *ConnectionApp) GetConnectionApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConnectionAppFilter, limit, offset int) (int, *[]model.ConnectionApp, error) {
	result := new([]model.ConnectionApp)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.ConnectionName) > 0 {
		query.Where("connection_name = ?", filter.ConnectionName)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
