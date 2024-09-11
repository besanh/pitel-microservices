package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatEmail interface {
		IRepo[model.ChatEmail]
		GetChatEmails(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatEmailFilter, limit, offset int) (total int, entries *[]model.ChatEmail, err error)
		GetChatEmailsCustom(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatEmailFilter, limit, offset int) (total int, entries *[]model.ChatEmailCustom, err error)
	}
	ChatEmail struct {
		Repo[model.ChatEmail]
	}
)

func NewChatEmail() IChatEmail {
	return &ChatEmail{}
}

func (repo *ChatEmail) GetChatEmails(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatEmailFilter, limit, offset int) (total int, entries *[]model.ChatEmail, err error) {
	entries = new([]model.ChatEmail)
	query := db.GetDB().NewSelect().
		Model(entries)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.OaId) > 0 {
		query.Where("oa_id = ?", filter.OaId)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status)
	}
	query.Order("created_at desc")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}

	return total, entries, nil
}

func (repo *ChatEmail) GetChatEmailsCustom(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatEmailFilter, limit, offset int) (total int, entries *[]model.ChatEmailCustom, err error) {
	entries = new([]model.ChatEmailCustom)
	query := db.GetDB().NewSelect().
		Model(entries).
		ColumnExpr("ce.*").
		ColumnExpr("tmp.*")
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	query2 := db.GetDB().NewSelect().TableExpr("chat_connection_app cca").
		ColumnExpr("connection_name").
		ColumnExpr("connection_type").
		ColumnExpr("oa_info").
		Where("cca.tenant_id = ?", filter.TenantId).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("cca.oa_info->'zalo'::text->0->>'oa_id' = ce.oa_id").
				WhereOr("cca.oa_info->'facebook'::text->0->>'oa_id' = ce.oa_id")
		})

	query.Join("LEFT JOIN LATERAL (?) tmp ON true", query2)

	if len(filter.OaId) > 0 {
		query.Where("ce.oa_id = ?", filter.OaId)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status)
	}
	query.Order("created_at desc")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}

	return total, entries, nil
}
