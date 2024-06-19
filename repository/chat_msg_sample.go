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
	IChatMsgSample interface {
		IRepo[model.ChatMsgSample]
		GetChatMsgSamples(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatMsgSampleFilter, limit, offset int) (int, *[]model.ChatMsgSampleView, error)
		GetChatMsgSampleById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatMsgSampleView, error)
	}

	ChatMsgSample struct {
		Repo[model.ChatMsgSample]
	}
)

var ChatMsgSampleRepo IChatMsgSample

func NewChatMsgSample() IChatMsgSample {
	return &ChatMsgSample{}
}

func (repo *ChatMsgSample) GetChatMsgSamples(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatMsgSampleFilter, limit, offset int) (int, *[]model.ChatMsgSampleView, error) {
	result := new([]model.ChatMsgSampleView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cms.*").
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name", "oa_info")
		})
	if len(filter.TenantId) > 0 {
		query.Where("cms.tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConnectionId) > 0 {
		query.Where("cms.connection_id = ?", filter.ConnectionId)
	}
	if len(filter.Channel) > 0 {
		query.Where("cms.channel = ?", filter.Channel)
	}
	if len(filter.OaId) > 0 {
		query.Where("connection_app.oa_info->cms.channel::text->0->>'oa_id' = ?", filter.OaId)
	}
	if len(filter.Keyword) > 0 {
		keyword := "%" + filter.Keyword + "%"
		query.Where("? ILIKE ?", bun.Ident("cms.keyword"), keyword).
			WhereOr("? ILIKE ?", bun.Ident("cms.theme"), keyword)
	}
	query.Order("cms.created_at desc")

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

func (repo *ChatMsgSample) GetChatMsgSampleById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatMsgSampleView, error) {
	result := new(model.ChatMsgSampleView)
	err := db.GetDB().NewSelect().
		Model(result).
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name")
		}).
		Where("cms.id = ?", id).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}
