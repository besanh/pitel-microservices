package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IChatConnectionPipeline interface {
		UpdateConnectionApp(ctx context.Context, tx bun.Tx, entity model.ChatConnectionApp) error
		InsertConnectionApp(ctx context.Context, tx bun.Tx, entity model.ChatConnectionApp) error
		DeleteConnectionQueue(ctx context.Context, tx bun.Tx, connectionId, queueId string) (err error)
		BeginTx(ctx context.Context, db sqlclient.ISqlClientConn, opts *sql.TxOptions) (bun.Tx, error)
		CommitTx(ctx context.Context, tx bun.Tx) error
	}
	ChatConnectionPipeline struct {
	}
)

var ChatConnectionPipelineRepo IChatConnectionPipeline

func NewConnectionPipeline() IChatConnectionPipeline {
	return &ChatConnectionPipeline{}
}

func (repo *ChatConnectionPipeline) DeleteConnectionQueue(ctx context.Context, tx bun.Tx, connectionId, queueId string) (err error) {
	query := tx.NewDelete().
		Model((*model.ConnectionQueue)(nil))
	if len(connectionId) > 0 {
		query.Where("connection_id = ?", connectionId)
	}
	if len(queueId) > 0 {
		query.Where("queue_id = ?", queueId)
	}
	_, err = query.Exec(ctx)
	return
}

func (repo *ChatConnectionPipeline) UpdateConnectionApp(ctx context.Context, tx bun.Tx, entity model.ChatConnectionApp) error {
	entity.UpdatedAt = time.Now()
	_, err := tx.NewUpdate().
		Model(&entity).
		Where("id = ?", entity.Id).
		Exec(ctx)
	return err
}

func (repo *ChatConnectionPipeline) InsertConnectionApp(ctx context.Context, tx bun.Tx, entity model.ChatConnectionApp) error {
	entity.CreatedAt = time.Now()
	_, err := tx.NewInsert().
		Model(&entity).
		Exec(ctx)
	return err
}

func (repo *ChatConnectionPipeline) BeginTx(ctx context.Context, db sqlclient.ISqlClientConn, opts *sql.TxOptions) (bun.Tx, error) {
	return db.GetDB().BeginTx(ctx, opts)
}

func (repo *ChatConnectionPipeline) CommitTx(ctx context.Context, tx bun.Tx) error {
	return tx.Commit()
}
