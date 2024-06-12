package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatCommand interface {
		IRepo[model.ChatCommand]
		GetChatCommands(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, []model.ChatCommandView, error)
	}

	ChatCommand struct {
		Repo[model.ChatCommand]
	}
)

var ChatCommandRepo IChatCommand

func NewChatCommand() IChatCommand {
	return &ChatCommand{}
}

func (repo *ChatCommand) GetChatCommands(ctx context.Context, db sqlclient.ISqlClientConn, limit, offset int) (int, []model.ChatCommandView, error) {
	result := make([]model.ChatCommandView, 0)
	query := db.GetDB().NewSelect().Model(result).
		Column("cc.*, cca.connection_name")
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
