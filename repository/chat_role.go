package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
)

type (
	IChatRole interface {
		IRepo[model.ChatRole]
		GetChatRoles(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatRoleFilter, limit, offset int) (total int, result *[]model.ChatRole, err error)
	}
	ChatRole struct {
		Repo[model.ChatRole]
	}
)

var ChatRoleRepo IChatRole

func NewChatRole() IChatRole {
	repo := &ChatRole{}
	return repo
}

func (repo *ChatRole) GetChatRoles(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatRoleFilter, limit, offset int) (total int, result *[]model.ChatRole, err error) {
	result = new([]model.ChatRole)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.RoleName) > 0 {
		query.Where("role_name = ?", filter.RoleName)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
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
