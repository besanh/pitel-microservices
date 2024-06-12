package repository

import (
	"context"
	"database/sql"
	"fmt"
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

	IChatPersonalization interface {
		IRepo[model.ChatPersonalization]
		InsertDefaultPersonalizationValue(ctx context.Context, db sqlclient.ISqlClientConn, str string) error
		GetPersonalizationValues(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, *[]model.ChatPersonalizationView, error)
	}

	ChatPersonalization struct {
		Repo[model.ChatPersonalization]
	}
)

var ChatMsgSampleRepo IChatMsgSample

func NewChatMsgSample() IChatMsgSample {
	return &ChatMsgSample{}
}

var ChatPersonalizationRepo IChatPersonalization

func NewChatPersonalization() IChatPersonalization {
	return &ChatPersonalization{}
}

func (repo *ChatMsgSample) GetChatMsgSamples(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, *[]model.ChatMsgSampleView, error) {
	result := new([]model.ChatMsgSampleView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cc.*").
		ColumnExpr("cca.connection_name")
	query.Join("LEFT JOIN chat_connection_app as cca").JoinOn("cc.id = cca.id")

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

func (repo *ChatPersonalization) InsertDefaultPersonalizationValue(ctx context.Context, db sqlclient.ISqlClientConn, str string) error {
	return repo.Insert(ctx, db, model.ChatPersonalization{
		Base:                 model.InitBase(),
		PersonalizationValue: fmt.Sprintf("{{%v}}", str),
	})
}

func (repo *ChatPersonalization) GetPersonalizationValues(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, *[]model.ChatPersonalizationView, error) {
	result := new([]model.ChatPersonalizationView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cp.personalization_value")

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
