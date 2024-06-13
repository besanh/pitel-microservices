package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
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
		Column("cst.*").
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name", "oa_info")
		})
	if len(filter.ScriptName) > 0 {
		query.Where("cst.script_name = ?", filter.ScriptName)
	}
	if len(filter.Channel) > 0 {
		query.Where("cst.channel = ?", filter.Channel)
	}
	if len(filter.OaId) > 0 {
		query.Where("connection_app.oa_info->cst.channel::text->0->>'oa_id' = ?", filter.OaId)
	}

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
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name")
		}).
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
