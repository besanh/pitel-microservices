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
		InsertConnectionQueue(ctx context.Context, tx bun.Tx, entity model.ConnectionQueue) error
		BeginTx(ctx context.Context, db sqlclient.ISqlClientConn, opts *sql.TxOptions) (bun.Tx, error)
	}
	ChatConnectionPipeline struct {
	}
)

var ChatConnectionPipelineRepo IChatConnectionPipeline

func NewConnectionPipeline() IChatConnectionPipeline {
	return &ChatConnectionPipeline{}
}

func (repo *ChatConnectionPipeline) InsertConnectionQueue(ctx context.Context, tx bun.Tx, entity model.ConnectionQueue) error {
	entity.SetCreatedAt(time.Now())
	_, err := tx.NewInsert().Model(entity).Exec(ctx)
	return err
}

func (repo *ChatConnectionPipeline) BeginTx(ctx context.Context, db sqlclient.ISqlClientConn, opts *sql.TxOptions) (bun.Tx, error) {
	return db.GetDB().BeginTx(ctx, opts)
}
