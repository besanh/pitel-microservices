package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatConnectionApp interface {
		GetChatConnectionApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionApp, error)
		GetChatConnectionAppCustom(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionAppView, error)
		GetById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (entity *model.ChatConnectionApp, err error)
		Insert(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatConnectionApp) (err error)
		Update(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatConnectionApp) (err error)
		Delete(ctx context.Context, db sqlclient.ISqlClientConn, id string) (err error)
		SelectByQuery(ctx context.Context, db sqlclient.ISqlClientConn, params []model.Param, limit int, offset int) (entries *[]model.ChatConnectionApp, total int, err error)
		BulkInsert(ctx context.Context, db sqlclient.ISqlClientConn, entities []model.ChatConnectionApp) (err error)
		UpdateSpecifColumnyById(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatConnectionApp) (err error)
	}
	ChatConnectionApp struct {
	}
)

var ChatConnectionAppRepo IChatConnectionApp

func NewConnectionApp() IChatConnectionApp {
	return &ChatConnectionApp{}
}

// Current: one connection having 1 element zalo/fb in 1 record
// TODO: one connection having many elements zalo/fb in 1 record
func (repo *ChatConnectionApp) GetChatConnectionApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionApp, error) {
	result := new([]model.ChatConnectionApp)
	query := db.GetDB().NewSelect().Model(result).
		Column("cca.*")
	if len(filter.TenantId) > 0 {
		query.Where("cca.tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConnectionName) > 0 {
		query.Where("connection_name = ?", filter.ConnectionName)
	}
	if len(filter.ConnectionType) > 0 && len(filter.OaId) > 0 {
		query.Where("oa_info->?::text->0->>'oa_id' = ?", filter.ConnectionType, filter.OaId)
	} else if len(filter.OaId) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			query.Where("oa_info->'zalo'::text->0->>'oa_id' = ?", filter.OaId).
				WhereOr("oa_info->'facebook'::text->0->>'oa_id' = ?", filter.OaId)

			return q
		})
	}

	if len(filter.ConnectionType) > 0 {
		query.Where("connection_type = ?", filter.ConnectionType)
	}

	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.ConnectionQueueId) > 0 {
		query.Where("connection_queue_id = ?", filter.ConnectionQueueId)
	}
	if len(filter.Status) > 0 {
		query.Where("status = ?", filter.Status)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("cca.created_at DESC")

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}

/**
* Use for connection with share form
 */
func (repo *ChatConnectionApp) GetChatConnectionAppCustom(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionAppView, error) {
	result := new([]model.ChatConnectionAppView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cca.*").
		ColumnExpr("tmp.*").
		Relation("ConnectionQueue")
	if len(filter.TenantId) > 0 {
		query.Where("cca.tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConnectionName) > 0 {
		query.Where("connection_name = ?", filter.ConnectionName)
	}
	if len(filter.ConnectionType) > 0 {
		query.Where("connection_type = ?", filter.ConnectionType)
		if len(filter.OaId) > 0 {
			query.Where("oa_info->?::text->0->>'oa_id' = ?", filter.ConnectionType, filter.OaId)
		}
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id = ?", filter.QueueId)
	}
	if len(filter.ConnectionQueueId) > 0 {
		query.Where("connection_queue_id = ?", filter.ConnectionQueueId)
	}
	if len(filter.Status) > 0 {
		query.Where("status = ?", filter.Status)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	query2 := db.GetDB().NewSelect().TableExpr("chat_share_info as csi").
		ColumnExpr("csi.id as share_form_uuid").
		ColumnExpr("csi.share_form as share_info_form").
		Where("csi.connection_id=cca.id")
	if len(filter.TenantId) > 0 {
		query2.Where("tenant_id = ?", filter.TenantId)
	}

	query.Join("LEFT JOIN LATERAL (?) AS tmp ON true", query2)
	query.Order("cca.created_at DESC")

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}

func (r *ChatConnectionApp) GetById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (entity *model.ChatConnectionApp, err error) {
	entity = new(model.ChatConnectionApp)
	query := db.GetDB().NewSelect().
		Model(entity).
		Relation("ConnectionQueue").
		Where("cca.id = ?", id).
		Limit(1)
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ChatConnectionApp) Insert(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatConnectionApp) (err error) {
	entity.CreatedAt = time.Now()
	_, err = db.GetDB().NewInsert().
		Model(&entity).
		Exec(ctx)
	return
}

func (r *ChatConnectionApp) Update(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatConnectionApp) (err error) {
	entity.UpdatedAt = time.Now()
	_, err = db.GetDB().NewUpdate().
		Model(&entity).
		Where("id = ?", entity.Id).
		Exec(ctx)
	return
}

func (r *ChatConnectionApp) Delete(ctx context.Context, db sqlclient.ISqlClientConn, id string) (err error) {
	_, err = db.GetDB().NewDelete().
		Model((*model.ChatConnectionApp)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return
}

func (r *ChatConnectionApp) SelectByQuery(ctx context.Context, db sqlclient.ISqlClientConn, params []model.Param, limit int, offset int) (entries *[]model.ChatConnectionApp, total int, err error) {
	entries = new([]model.ChatConnectionApp)
	query := db.GetDB().NewSelect().
		Model(entries).
		Limit(limit).
		Offset(offset)
	for _, param := range params {
		qb := param.BuildQuery()
		if len(qb) > 0 {
			query.Where(qb, param.Value)
		}
	}
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return entries, total, nil
}

func (r *ChatConnectionApp) BulkInsert(ctx context.Context, db sqlclient.ISqlClientConn, entities []model.ChatConnectionApp) (err error) {
	_, err = db.GetDB().NewInsert().
		Model(&entities).
		Exec(ctx)
	return
}

func (r *ChatConnectionApp) UpdateSpecifColumnyById(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatConnectionApp) (err error) {
	entity.UpdatedAt = time.Now()
	_, err = db.GetDB().NewUpdate().
		Model(&entity).
		Column("updated_at", "connection_name", "oa_info", "app_id", "status").
		Where("id = ?", entity.Id).
		Exec(ctx)
	return
}
