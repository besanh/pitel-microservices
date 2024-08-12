package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository/sql_builder"
)

type (
	IIBKBusinessUnit interface {
		IRepo[model.IBKBusinessUnit]
		SelectInfoByQuery(ctx context.Context, conditions []sql_builder.QueryCondition, limit int, offset int) (entries *[]model.IBKBusinessUnitInfo, total int, err error)
		FindInfoById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entry *model.IBKBusinessUnitInfo, err error)
	}
	IBKBusinessUnit struct {
		Repo[model.IBKBusinessUnit]
	}
)

var IBKBusinessUnitRepo IIBKBusinessUnit

func NewIBKBusinessUnit(conn sqlclient.ISqlClientConn) IIBKBusinessUnit {
	repo := &IBKBusinessUnit{
		Repo[model.IBKBusinessUnit]{
			Conn: conn,
		},
	}
	return repo
}

func (repo *IBKBusinessUnit) InitTable(ctx context.Context) {
	if err := CreateTable(ctx, repo.Conn, (*model.IBKBusinessUnit)(nil)); err != nil {
		log.Error(err)
	}
}

func (repo *IBKBusinessUnit) SelectInfoByQuery(ctx context.Context, conditions []sql_builder.QueryCondition, limit int, offset int) (entries *[]model.IBKBusinessUnitInfo, total int, err error) {
	entries = new([]model.IBKBusinessUnitInfo)
	query := repo.Conn.GetDB().NewSelect().
		Model(entries).
		Limit(limit).
		Offset(offset).
		ColumnExpr("bu.*").
		ColumnExpr("(?) as total_user", repo.Conn.GetDB().NewSelect().TableExpr("ibk_users as u").
			ColumnExpr("count(u.id) as total_user"),
		)
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}

	query.OrderExpr("bu.created_at DESC")
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return entries, total, nil
}

func (repo *IBKBusinessUnit) FindInfoById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entry *model.IBKBusinessUnitInfo, err error) {
	entry = new(model.IBKBusinessUnitInfo)
	query := repo.Conn.GetDB().NewSelect().
		Model(entry).
		ColumnExpr("bu.*").
		ColumnExpr("(?) as total_user", repo.Conn.GetDB().NewSelect().TableExpr("ibk_users as u").
			ColumnExpr("count(u.id) as total_user"),
		).
		Limit(1).
		Where("bu.id = ?", id)
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}
