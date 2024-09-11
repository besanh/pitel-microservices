package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
)

type (
	IChatUser interface {
		IRepo[model.ChatUser]
		GetChatUsers(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatUserFilter, limit, offset int) (total int, result *[]model.ChatUser, err error)
		GetChatUserByFilter(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatUserFilter) (result *model.ChatUser, err error)
	}
	ChatUser struct {
		Repo[model.ChatUser]
	}
)

var ChatUserRepo IChatUser

func NewChatUser() IChatUser {
	return &ChatUser{}
}

func (repo *ChatUser) GetChatUsers(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatUserFilter, limit, offset int) (total int, result *[]model.ChatUser, err error) {
	result = new([]model.ChatUser)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.Username) > 0 {
		query.Where("username = ?", filter.Username)
	}
	if len(filter.Level) > 0 {
		query.Where("level = ?", filter.Level)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if len(filter.Fullname) > 0 {
		query.Where("fullname = ?", filter.Fullname)
	}
	if len(filter.RoleId) > 0 {
		query.Where("role_id = ?", filter.RoleId)
	}
	query.Order("created_at desc")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}

func (repo *ChatUser) GetChatUserByFilter(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatUserFilter) (result *model.ChatUser, err error) {
	result = new(model.ChatUser)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.Username) > 0 {
		query.Where("username = ?", filter.Username)
	}
	if len(filter.Level) > 0 {
		query.Where("level = ?", filter.Level)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if len(filter.Fullname) > 0 {
		query.Where("fullname = ?", filter.Fullname)
	}
	if len(filter.RoleId) > 0 {
		query.Where("role_id = ?", filter.RoleId)
	}
	err = query.Limit(1).Scan(ctx)
	return result, err
}
