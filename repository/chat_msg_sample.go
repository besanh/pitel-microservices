package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatMsgSample interface {
		IRepo[model.ChatMsgSample]
		GetChatMsgSamples(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, *[]model.ChatMsgSampleView, error)
	}

	ChatMsgSample struct {
		Repo[model.ChatMsgSample]
	}
)

var ChatMsgSampleRepo IChatMsgSample

func NewChatMsgSample() IChatMsgSample {
	return &ChatMsgSample{}
}

func (repo *ChatMsgSample) GetChatMsgSamples(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, *[]model.ChatMsgSampleView, error) {
	result := new([]model.ChatMsgSampleView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cms.*").
		ColumnExpr("connection_app.connection_name").
		Relation("ConnectionApp")

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
