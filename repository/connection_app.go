package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatConnectionApp interface {
		IRepo[model.ChatConnectionApp]
		GetChatConnectionApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionApp, error)
	}
	ChatConnectionApp struct {
		Repo[model.ChatConnectionApp]
	}
)

var ChatConnectionAppRepo IChatConnectionApp

func NewConnectionApp() IChatConnectionApp {
	return &ChatConnectionApp{}
}

func (repo *ChatConnectionApp) GetChatConnectionApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionApp, error) {
	result := new([]model.ChatConnectionApp)
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
