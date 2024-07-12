package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatApp interface {
		IRepo[model.ChatApp]
		GetChatApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAppFilter, limit, offset int) (int, *[]model.ChatApp, error)
		GetChatAppById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatApp, error)
		InsertChatApp(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatApp, systems []model.ChatAppIntegrateSystem) error
		UpdateChatAppById(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatApp, systems []model.ChatAppIntegrateSystem) error
		DeleteChatAppById(ctx context.Context, db sqlclient.ISqlClientConn, id string) error
	}
	ChatApp struct {
		Repo[model.ChatApp]
	}
)

var ChatAppRepo IChatApp

func NewChatApp() IChatApp {
	return &ChatApp{}
}

func (s *ChatApp) GetChatApp(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAppFilter, limit, offset int) (int, *[]model.ChatApp, error) {
	result := new([]model.ChatApp)
	query := db.GetDB().NewSelect().Model(result).
		Relation("Systems")
	if len(filter.AppName) > 0 {
		query.Where("app_name = ?", filter.AppName)
	}
	if len(filter.Status) > 0 {
		query.Where("status = ?", filter.Status)
	}
	if len(filter.AppType) > 0 {
		if len(filter.Status) > 0 {
			query.Where("info_app :: jsonb -> ? ->> 'status' = 'active'", filter.AppType)
		} else {
			query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("info_app :: jsonb -> ? ->> 'status' = 'active'", filter.AppType).
					WhereOr("info_app :: jsonb -> ? ->> 'status' = 'deactive'", filter.AppType)
			})
		}
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at desc")
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}

	return total, result, nil
}

func (s *ChatApp) GetChatAppById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatApp, error) {
	result := new(model.ChatApp)
	err := db.GetDB().NewSelect().Model(result).
		Relation("Systems").
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

func (s *ChatApp) InsertChatApp(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatApp, systems []model.ChatAppIntegrateSystem) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.NewInsert().Model(&data).Exec(ctx); err != nil {
		return err
	}
	if len(systems) > 0 {
		if _, err = tx.NewInsert().Model(&systems).Exec(ctx); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *ChatApp) UpdateChatAppById(ctx context.Context, db sqlclient.ISqlClientConn, data model.ChatApp, systems []model.ChatAppIntegrateSystem) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewUpdate().Model(&data).
		Where("id = ?", data.GetId()).
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewDelete().Model((*model.ChatAppIntegrateSystem)(nil)).
		Where("chat_app_id = ?", data.GetId()).
		Exec(ctx)
	if err != nil {
		return err
	}

	if len(systems) > 0 {
		if _, err = tx.NewInsert().Model(&systems).Exec(ctx); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *ChatApp) DeleteChatAppById(ctx context.Context, db sqlclient.ISqlClientConn, id string) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewDelete().Model((*model.ChatAppIntegrateSystem)(nil)).
		Where("chat_app_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewDelete().Model((*model.ChatApp)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
