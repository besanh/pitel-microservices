package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IAllocateUser interface {
		IRepo[model.AllocateUser]
		GetAllocateUsers(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AllocateUserFilter, limit, offset int) (int, *[]model.AllocateUser, error)
		DeleteAllocateUsers(ctx context.Context, db sqlclient.ISqlClientConn, userAllocates []model.AllocateUser) error
	}
	UserAllocate struct {
		Repo[model.AllocateUser]
	}
)

var AllocateUserRepo IAllocateUser

func NewAllocateUser() IAllocateUser {
	return &UserAllocate{}
}

func (repo *UserAllocate) GetAllocateUsers(ctx context.Context, db sqlclient.ISqlClientConn, filter model.AllocateUserFilter, limit, offset int) (int, *[]model.AllocateUser, error) {
	result := new([]model.AllocateUser)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.AppId) > 0 {
		query.Where("app_id = ?", filter.AppId)
	}
	if len(filter.OaId) > 0 {
		query.Where("oa_id = ?", filter.OaId)
	}
	if len(filter.UserId) > 0 {
		query.Where("user_id IN (?)", bun.In(filter.UserId))
	}
	if len(filter.QueueId) > 0 {
		query.Where("queue_id IN (?)", bun.In(filter.QueueId))
	}
	if len(filter.ConversationId) > 0 {
		query.Where("conversation_id = ?", filter.ConversationId)
	}
	if len(filter.ExternalConversationId) > 0 {
		query.Where("external_conversation_id = ?", filter.ExternalConversationId)
	}
	if len(filter.MainAllocate) > 0 {
		query.Where("main_allocate = ?", filter.MainAllocate)
	}
	query.Order("created_at desc")
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

func (repo *UserAllocate) DeleteAllocateUsers(ctx context.Context, db sqlclient.ISqlClientConn, userAllocates []model.AllocateUser) error {
	_, err := db.GetDB().NewDelete().
		Model(&userAllocates).
		WherePK().
		Exec(ctx)

	return err
}
