package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository/sql_builder"
	"github.com/uptrace/bun"
)

type IRepo[T model.Model] interface {
	GetById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (*T, error)
	GetByQuery(ctx context.Context, conditions ...sql_builder.QueryCondition) (entry *T, err error)
	SelectByQuery(ctx context.Context, conditions []sql_builder.QueryCondition, limit int, offset int, orderBy string) (total int, entries []*T, err error)

	CreateTable(ctx context.Context) (err error)

	Insert(ctx context.Context, entry T) error
	Update(ctx context.Context, entry T) error
	Delete(ctx context.Context, id ...string) error

	UpdateByMap(ctx context.Context, id string, data map[string]interface{}) (err error)
	InsertIfNotExists(ctx context.Context, entry T) (err error)
	BulkDelete(ctx context.Context, ids []string) (err error)

	// transaction
	SelectByQueryWithTx(ctx context.Context, tx bun.Tx, conditions []sql_builder.QueryCondition, limit int, offset int, orderBy string) (entries []*T, total int, err error)
	InsertWithTx(ctx context.Context, tx bun.Tx, entries ...T) (err error)
	UpdateWithTx(ctx context.Context, tx bun.Tx, entry T) (err error)
	InsertTx(ctx context.Context, entry T) (err error)
	UpdateTx(ctx context.Context, entry T) (err error)
	InsertOrUpdateTx(ctx context.Context, entry T) (err error)
	CreateTx(ctx context.Context) (tx bun.Tx, err error)
	DeleteWithTx(ctx context.Context, tx bun.Tx, ids ...string) (err error)
	UpdateByMapWithTx(ctx context.Context, tx bun.Tx, id string, data map[string]interface{}) (err error)

	InitIndexes(ctx context.Context)
	AlterColumns(ctx context.Context)
	InitTable(ctx context.Context)

	CreateTxWithDB(ctx context.Context, conn sqlclient.ISqlClientConn) (tx bun.Tx, err error)
	GetByIdWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, id string, conditions ...sql_builder.QueryCondition) (*T, error)
	SelectByQueryWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, conditions []sql_builder.QueryCondition, limit int, offset int, orderBy string) (entries []*T, total int, err error)
	GetByQueryWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, conditions ...sql_builder.QueryCondition) (entry *T, err error)

	InsertWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, entries ...T) error
	UpdateWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, entries ...T) error
	DeleteWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, ids ...string) error
	UpdateByMapWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, id string, data map[string]interface{}) (err error)
	InsertIfNotExistsWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, entries ...T) (err error)

	InitTableWithDB(ctx context.Context, conn sqlclient.ISqlClientConn)
	AlterColumnsWithDB(ctx context.Context, conn sqlclient.ISqlClientConn)
	InitIndexesWithDB(ctx context.Context, conn sqlclient.ISqlClientConn)
}

type Repo[T model.Model] struct {
	Conn sqlclient.ISqlClientConn
}

func NewRepo[T model.Model](conn sqlclient.ISqlClientConn) IRepo[T] {
	return &Repo[T]{
		Conn: conn,
	}
}

func CreateTable(ctx context.Context, conn sqlclient.ISqlClientConn, entry any) (err error) {
	_, err = conn.GetDB().NewCreateTable().Model(entry).
		IfNotExists().
		Exec(ctx)
	return
}

func (repo *Repo[T]) CreateTable(ctx context.Context) (err error) {
	query := repo.Conn.GetDB().NewCreateTable().Model((*T)(nil)).
		IfNotExists()
	_, err = query.
		Exec(ctx)
	return
}

func (repo *Repo[T]) GetById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entry *T, err error) {
	entry = new(T)
	query := repo.Conn.GetDB().NewSelect().
		Model(entry).
		Where("id = ?", id).
		Limit(1)
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

func (repo *Repo[T]) Insert(ctx context.Context, entry T) (err error) {
	entry.SetCreatedAt(time.Now())
	_, err = repo.Conn.GetDB().NewInsert().
		Model(&entry).
		Exec(ctx)
	return
}

func (repo *Repo[T]) Update(ctx context.Context, entry T) (err error) {
	entry.SetUpdatedAt(time.Now())
	_, err = repo.Conn.GetDB().NewUpdate().
		Model(&entry).
		Where("id = ?", entry.GetId()).
		Exec(ctx)
	return
}

func (repo *Repo[T]) Delete(ctx context.Context, id ...string) (err error) {
	_, err = repo.Conn.GetDB().NewDelete().
		Model((*T)(nil)).
		Where("id IN (?)", bun.In(id)).
		Exec(ctx)
	return
}

func (repo *Repo[T]) SelectByQuery(ctx context.Context, conditions []sql_builder.QueryCondition, limit int, offset int, orderBy string) (total int, entries []*T, err error) {
	entries = make([]*T, 0)
	query := repo.Conn.GetDB().NewSelect().
		Model(&entries).
		Offset(offset)
	if limit > 0 {
		query.Limit(limit)
	}
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}
	if len(orderBy) > 0 {
		query.Order(orderBy)
	}
	log.Infof("query: %v", query.String())
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}
	return total, entries, nil
}

func (repo *Repo[T]) GetByQuery(ctx context.Context, conditions ...sql_builder.QueryCondition) (entry *T, err error) {
	entry = new(T)
	query := repo.Conn.GetDB().NewSelect().
		Model(entry).
		Limit(1)
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}
	log.Infof("query: %v", query.String())
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (repo *Repo[T]) BulkDelete(ctx context.Context, ids []string) (err error) {
	return repo.Conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, id := range ids {
			if _, err = tx.NewDelete().
				Model((*T)(nil)).
				Where("id = ?", id).
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Repo[T]) UpdateByMap(ctx context.Context, id string, data map[string]interface{}) (err error) {
	query := repo.Conn.GetDB().NewUpdate().
		Model((*T)(nil)).
		Where("id = ?", id)
	for k, v := range data {
		if k == "id" {
			continue
		}
		query.Set(fmt.Sprintf("%s = ?", k), v)
	}
	query.Set("updated_at = ?", time.Now())
	_, err = query.
		Exec(ctx)
	return
}

func (repo *Repo[T]) InsertIfNotExists(ctx context.Context, entry T) (err error) {
	_, err = repo.Conn.GetDB().NewInsert().
		Model(&entry).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	return
}

func (repo *Repo[T]) SelectByQueryWithTx(ctx context.Context, tx bun.Tx, conditions []sql_builder.QueryCondition, limit int, offset int, orderBy string) (entries []*T, total int, err error) {
	entries = make([]*T, 0)
	query := tx.NewSelect().
		Model(&entries).
		Offset(offset)
	if limit > 0 {
		query.Limit(limit)
	}
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}
	if len(orderBy) > 0 {
		query.Order(orderBy)
	}
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return entries, total, nil
}

func (repo *Repo[T]) InsertTx(ctx context.Context, entry T) (err error) {
	return repo.Conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewInsert().
			Model(&entry).
			Exec(ctx)
		return err
	})
}

func (repo *Repo[T]) UpdateTx(ctx context.Context, entry T) (err error) {
	return repo.Conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		entry.SetUpdatedAt(time.Now())
		query := tx.NewUpdate().
			Model(&entry).
			Where("id = ?", entry.GetId())
		_, err = query.
			Exec(ctx)
		return err
	})
}

func (repo *Repo[T]) InsertOrUpdateTx(ctx context.Context, entry T) (err error) {
	return repo.Conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewInsert().
			Model(&entry).
			On("CONFLICT (id) DO UPDATE").
			Exec(ctx)
		return err
	})
}

func (repo *Repo[T]) CreateTx(ctx context.Context) (tx bun.Tx, err error) {
	tx, err = repo.Conn.GetDB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}
	return
}

func (repo *Repo[T]) InsertWithTx(ctx context.Context, tx bun.Tx, entries ...T) (err error) {
	for _, entry := range entries {
		_, err = tx.NewInsert().
			Model(&entry).
			Exec(ctx)
	}
	return
}

func (repo *Repo[T]) UpdateWithTx(ctx context.Context, tx bun.Tx, entry T) (err error) {
	entry.SetUpdatedAt(time.Now())
	query := tx.NewUpdate().
		Model(&entry).
		Where("id = ?", entry.GetId())
	_, err = query.
		Exec(ctx)
	return
}

func (repo *Repo[T]) DeleteWithTx(ctx context.Context, tx bun.Tx, ids ...string) (err error) {
	for _, id := range ids {
		_, err = tx.NewDelete().
			Model((*T)(nil)).
			Where("id = ?", id).
			Exec(ctx)
		if err != nil {
			return
		}
	}
	return
}

func (repo *Repo[T]) UpdateByMapWithTx(ctx context.Context, tx bun.Tx, id string, data map[string]interface{}) (err error) {
	query := tx.NewUpdate().
		Model((*T)(nil)).
		Where("id = ?", id)
	for k, v := range data {
		if k == "id" {
			continue
		}
		query.Set(fmt.Sprintf("%s = ?", k), v)
	}
	query.Set("updated_at = ?", time.Now())
	_, err = query.
		Exec(ctx)
	return
}

func (repo *Repo[T]) InitIndexes(ctx context.Context) {}

func (repo *Repo[T]) AlterColumns(ctx context.Context) {}

func (repo *Repo[T]) InitTable(ctx context.Context) {}

func (repo *Repo[T]) InitTableWithDB(ctx context.Context, conn sqlclient.ISqlClientConn) {}

func (repo *Repo[T]) AlterColumnsWithDB(ctx context.Context, conn sqlclient.ISqlClientConn) {}

func (repo *Repo[T]) InitIndexesWithDB(ctx context.Context, conn sqlclient.ISqlClientConn) {}

func (repo *Repo[T]) CreateTxWithDB(ctx context.Context, conn sqlclient.ISqlClientConn) (tx bun.Tx, err error) {
	tx, err = conn.GetDB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}
	return
}

func (repo *Repo[T]) GetByIdWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, id string, conditions ...sql_builder.QueryCondition) (entry *T, err error) {
	entry = new(T)
	query := conn.GetDB().NewSelect().
		Model(entry).
		Where("id = ?", id).
		Limit(1)
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

func (repo *Repo[T]) GetByQueryWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, conditions ...sql_builder.QueryCondition) (entry *T, err error) {
	entry = new(T)
	query := conn.GetDB().NewSelect().
		Model(entry).
		Limit(1)
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}
	log.Infof("query: %v", query.String())
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}
func (repo *Repo[T]) SelectByQueryWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, conditions []sql_builder.QueryCondition, limit int, offset int, orderBy string) (entries []*T, total int, err error) {
	entries = make([]*T, 0)
	query := conn.GetDB().NewSelect().
		Model(&entries).
		Offset(offset)
	if limit > 0 {
		query.Limit(limit)
	}
	for _, condition := range conditions {
		qb := sql_builder.BuildQuery(condition)
		if len(qb) > 0 {
			query.Where(qb)
		}
	}
	if len(orderBy) > 0 {
		query.Order(orderBy)
	}
	log.Debugf("query: %v", query.String())
	total, err = query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return entries, total, nil
}

func (repo *Repo[T]) InsertWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, entries ...T) (err error) {
	return conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, entry := range entries {
			_, err = tx.NewInsert().
				Model(&entry).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Repo[T]) InsertIfNotExistsWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, entries ...T) (err error) {
	return conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, entry := range entries {
			_, err = tx.NewInsert().
				Model(&entry).
				On("CONFLICT (id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Repo[T]) UpdateWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, entries ...T) (err error) {
	return conn.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, entry := range entries {
			entry.SetUpdatedAt(time.Now())

			_, err = tx.NewUpdate().
				Model(&entry).
				Where("id = ?", entry.GetId()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Repo[T]) DeleteWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, ids ...string) (err error) {
	for _, id := range ids {
		_, err = conn.GetDB().NewDelete().
			Model((*T)(nil)).
			Where("id = ?", id).
			Exec(ctx)
		return
	}
	return
}

func (repo *Repo[T]) UpdateByMapWithDB(ctx context.Context, conn sqlclient.ISqlClientConn, id string, data map[string]interface{}) (err error) {
	query := conn.GetDB().NewUpdate().
		Model((*T)(nil)).
		Where("id = ?", id)
	for k, v := range data {
		if k == "id" {
			continue
		}
		query.Set(fmt.Sprintf("%s = ?", k), v)
	}
	query.Set("updated_at = ?", time.Now())
	_, err = query.
		Exec(ctx)
	return
}
