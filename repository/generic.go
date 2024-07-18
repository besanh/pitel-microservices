package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type IRepo[T model.Model] interface {
	GetById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*T, error)
	Insert(ctx context.Context, db sqlclient.ISqlClientConn, entity T) error
	Update(ctx context.Context, db sqlclient.ISqlClientConn, entity T) error
	Delete(ctx context.Context, db sqlclient.ISqlClientConn, id string) error
	CreateTable(ctx context.Context, db sqlclient.ISqlClientConn) (err error)
	SelectByQuery(ctx context.Context, db sqlclient.ISqlClientConn, params []model.Param, limit int, offset int) (entries *[]T, total int, err error)
	BulkInsert(ctx context.Context, db sqlclient.ISqlClientConn, entities []T) error
	TxInsert(ctx context.Context, tx bun.Tx, entity T) error
	TxBulkInsert(ctx context.Context, tx bun.Tx, entities []T) error
	TxUpdate(ctx context.Context, tx bun.Tx, entity T) error
	TxBulkUpdate(ctx context.Context, tx bun.Tx, entities []T) error
	TxDelete(ctx context.Context, tx bun.Tx, entity T) error
	TxBulkDelete(ctx context.Context, tx bun.Tx, entities []T) error
}

type Repo[T model.Model] struct {
}

func NewRepo[T model.Model]() IRepo[T] {
	return &Repo[T]{}
}

func (r *Repo[T]) CreateTable(ctx context.Context, db sqlclient.ISqlClientConn) (err error) {
	query := db.GetDB().NewCreateTable().Model((*T)(nil)).
		IfNotExists()
	// value, _ := query.AppendQuery(schema.NewFormatter(query.Dialect()), nil)
	// queryStr := string(value)
	// log.Infof("query: %v", queryStr)
	_, err = query.
		Exec(ctx)
	return
}

func (r *Repo[T]) GetById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (entity *T, err error) {
	entity = new(T)
	err = db.GetDB().NewSelect().
		Model(entity).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *Repo[T]) Insert(ctx context.Context, db sqlclient.ISqlClientConn, entity T) (err error) {
	entity.SetCreatedAt(time.Now())
	_, err = db.GetDB().NewInsert().
		Model(&entity).
		Exec(ctx)
	return
}

func (r *Repo[T]) Update(ctx context.Context, db sqlclient.ISqlClientConn, entity T) (err error) {
	entity.SetUpdatedAt(time.Now())
	_, err = db.GetDB().NewUpdate().
		Model(&entity).
		Where("id = ?", entity.GetId()).
		Exec(ctx)
	return
}

func (r *Repo[T]) Delete(ctx context.Context, db sqlclient.ISqlClientConn, id string) (err error) {
	_, err = db.GetDB().NewDelete().
		Model((*T)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return
}

func (r *Repo[T]) SelectByQuery(ctx context.Context, db sqlclient.ISqlClientConn, params []model.Param, limit int, offset int) (entries *[]T, total int, err error) {
	entries = new([]T)
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

func (r *Repo[T]) BulkInsert(ctx context.Context, db sqlclient.ISqlClientConn, entities []T) (err error) {
	_, err = db.GetDB().NewInsert().
		Model(&entities).
		Exec(ctx)
	return
}

func (r *Repo[T]) TxInsert(ctx context.Context, tx bun.Tx, entity T) (err error) {
	entity.SetCreatedAt(time.Now())
	_, err = tx.NewInsert().
		Model(&entity).
		Exec(ctx)
	return
}

func (r *Repo[T]) TxBulkInsert(ctx context.Context, tx bun.Tx, entities []T) (err error) {
	_, err = tx.NewInsert().
		Model(&entities).
		Exec(ctx)
	return
}

func (r *Repo[T]) TxUpdate(ctx context.Context, tx bun.Tx, entity T) (err error) {
	entity.SetUpdatedAt(time.Now())
	_, err = tx.NewUpdate().
		Model(&entity).
		WherePK().
		Exec(ctx)
	return
}

func (r *Repo[T]) TxBulkUpdate(ctx context.Context, tx bun.Tx, entities []T) (err error) {
	_, err = tx.NewUpdate().
		Model(&entities).
		Bulk().
		Exec(ctx)
	return
}

func (r *Repo[T]) TxDelete(ctx context.Context, tx bun.Tx, entity T) (err error) {
	_, err = tx.NewDelete().
		Model(&entity).
		WherePK().
		Exec(ctx)
	return
}

func (r *Repo[T]) TxBulkDelete(ctx context.Context, tx bun.Tx, entities []T) (err error) {
	_, err = tx.NewDelete().
		Model(&entities).
		Where("id IN (?)", bun.In(r.getIds(entities))).
		Exec(ctx)
	return
}

func (r *Repo[T]) getIds(entities []T) (result []string) {
	result = make([]string, len(entities))
	for i, entity := range entities {
		result[i] = entity.GetId()
	}
	return
}
