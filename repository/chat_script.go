package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatScript interface {
		IRepo[model.ChatScript]
		GetChatScripts(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatScriptFilter, limit, offset int) (int, *[]model.ChatScriptView, error)
		GetChatScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatScriptView, error)
	}

	ChatScript struct {
		Repo[model.ChatScript]
	}
)

var ChatScriptRepo IChatScript

func NewChatScript() IChatScript {
	return &ChatScript{}
}

func (repo *ChatScript) GetChatScripts(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatScriptFilter, limit, offset int) (int, *[]model.ChatScriptView, error) {
	result := new([]model.ChatScriptView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cst.*")
	if len(filter.ScriptName) > 0 {
		query.Where("cst.script_name ILIKE ?", "%"+filter.ScriptName+"%")
	}
	if len(filter.TenantId) > 0 {
		query.Where("cst.tenant_id = ?", filter.TenantId)
	}
	if len(filter.Channel) > 0 {
		query.Where("cst.channel = ?", filter.Channel)
	}
	if filter.Status.Valid {
		query.Where("cst.status = ?", filter.Status.Bool)
	}
	query.Order("cst.created_at desc")

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, result, nil
	} else if err != nil {
		return 0, nil, err
	}
	return total, result, nil
}

func (repo *ChatScript) GetChatScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatScriptView, error) {
	result := new(model.ChatScriptView)
	err := db.GetDB().NewSelect().
		Model(result).
		Where("cst.id = ?", id).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}
