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
	IChatAutoScript interface {
		IRepo[model.ChatAutoScript]
		GetChatAutoScripts(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAutoScriptFilter, limit, offset int) (int, *[]model.ChatAutoScriptView, error)
		GetChatAutoScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatAutoScriptView, error)
	}

	ChatAutoScript struct {
		Repo[model.ChatAutoScript]
	}
)

var ChatAutoScriptRepo IChatAutoScript

func NewChatAutoScript() IChatAutoScript {
	return &ChatAutoScript{}
}

func (repo *ChatAutoScript) GetChatAutoScripts(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAutoScriptFilter, limit, offset int) (int, *[]model.ChatAutoScriptView, error) {
	result := new([]model.ChatAutoScriptView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cas.*").
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name", "oa_info")
		})
	if len(filter.ScriptName) > 0 {
		query.Where("cas.script_name = ?", filter.ScriptName)
	}
	if len(filter.Channel) > 0 {
		query.Where("cas.channel = ?", filter.Channel)
	}
	if len(filter.OaId) > 0 {
		query.Where("connection_app.oa_info->cas.channel::text->0->>'oa_id' = ?", filter.OaId)
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

func (repo *ChatAutoScript) GetChatAutoScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatAutoScriptView, error) {
	result := new(model.ChatAutoScriptView)
	err := db.GetDB().NewSelect().
		Model(result).
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name")
		}).
		Where("cas.id = ?", id).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}
