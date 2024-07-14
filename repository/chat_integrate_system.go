package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatIntegrateSystem interface {
		IRepo[model.ChatIntegrateSystem]
		GetIntegrateSystems(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error)
		GetIntegrateSystemById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatIntegrateSystem, error)
		InsertIntegrateSystem(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatIntegrateSystem, chatApps []model.ChatAppIntegrateSystem) error
		UpdateIntegrateSystemById(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatIntegrateSystem, chatApps []model.ChatAppIntegrateSystem) error
		DeleteIntegrateSystemById(ctx context.Context, db sqlclient.ISqlClientConn, id string) error
	}
	ChatIntegrateSystem struct {
		Repo[model.ChatIntegrateSystem]
	}
)

var ChatIntegrateSystemRepo IChatIntegrateSystem

func NewChatIntegrateSystem() IChatIntegrateSystem {
	return &ChatIntegrateSystem{}
}

func (repo *ChatIntegrateSystem) GetIntegrateSystems(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error) {
	result = new([]model.ChatIntegrateSystem)
	query := db.GetDB().NewSelect().Model(result).
		Relation("Vendor", func(q *bun.SelectQuery) *bun.SelectQuery {
			if len(filter.VendorName) > 0 {
				return q.Where("vendor_name = ?", filter.VendorName)
			}
			return q
		}).
		Relation("ChatApps")
	if len(filter.SystemName) > 0 {
		query.Where("system_name = ?", filter.SystemName)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status)
	}
	if len(filter.SystemId) > 0 {
		query.Where("system_id = ?", filter.SystemId)
	}
	query.Order("created_at DESC")

	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}

func (repo *ChatIntegrateSystem) GetIntegrateSystemById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatIntegrateSystem, error) {
	result := new(model.ChatIntegrateSystem)
	err := db.GetDB().NewSelect().Model(result).
		Relation("ChatApps").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *ChatIntegrateSystem) InsertIntegrateSystem(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatIntegrateSystem, chatApps []model.ChatAppIntegrateSystem) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	data.SetCreatedAt(time.Now())
	if _, err = tx.NewInsert().Model(&data).Exec(ctx); err != nil {
		return err
	}
	if len(chatApps) > 0 {
		if _, err = tx.NewInsert().Model(&chatApps).Exec(ctx); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (repo *ChatIntegrateSystem) UpdateIntegrateSystemById(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatIntegrateSystem, chatApps []model.ChatAppIntegrateSystem) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	data.SetUpdatedAt(time.Now())
	_, err = tx.NewUpdate().Model(&data).
		Where("id = ?", data.GetId()).
		Exec(ctx)
	if err != nil {
		return err
	}

	if len(chatApps) > 0 {
		_, err = tx.NewDelete().Model((*model.ChatAppIntegrateSystem)(nil)).
			Where("chat_integrate_system_id = ?", data.GetId()).
			Exec(ctx)
		if err != nil {
			return err
		}
		if _, err = tx.NewInsert().Model(&chatApps).Exec(ctx); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (repo *ChatIntegrateSystem) DeleteIntegrateSystemById(ctx context.Context, db sqlclient.ISqlClientConn, id string) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewDelete().Model((*model.ChatAppIntegrateSystem)(nil)).
		Where("chat_integrate_system_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewDelete().Model((*model.ChatIntegrateSystem)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
