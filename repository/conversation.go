package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IConversation interface {
		IRepo[model.Conversation]
		GetConversations(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConversationFilter, limit, offset int) (int, *[]model.Conversation, error)
	}
	Conversation struct {
		Repo[model.Conversation]
	}
)

var ConversationRepo IConversation

func NewConversation() IConversation {
	return &Conversation{}
}

func (rpeo *Conversation) GetConversations(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConversationFilter, limit, offset int) (int, *[]model.Conversation, error) {
	result := new([]model.Conversation)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.UserIdByApp) > 0 {
		query.Where("user_id_by_app = ?", filter.UserIdByApp)
	}
	if len(filter.Username) > 0 {
		query.Where("username = ?", filter.Username)
	}
	if len(filter.PhoneNumber) > 0 {
		query.Where("phone_number = ?", filter.PhoneNumber)
	}
	if len(filter.Email) > 0 {
		query.Where("email = ?", filter.Email)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}

	return total, result, nil
}
