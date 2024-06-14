package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatLabel interface {
		IRepo[model.ChatLabel]
		GetChatLabels(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatLabelFilter, limit, offset int) (total int, result *[]model.ChatLabel, err error)
	}
	ChatLabel struct {
		Repo[model.ChatLabel]
	}
)

var ChatLabelRepo IChatLabel

func NewChatLabel() IChatLabel {
	repo := &ChatLabel{}
	return repo
}

func (repo *ChatLabel) GetChatLabels(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatLabelFilter, limit, offset int) (total int, result *[]model.ChatLabel, err error) {
	result = new([]model.ChatLabel)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.LabelName) > 0 {
		query.Where("Label_name = ?", filter.LabelName)
	}
	if len(filter.LabelColor) > 0 {
		query.Where("Label_color = ?", filter.LabelColor)
	}
	if filter.LabelStatus.Valid {
		query.Where("Label_status = ?", filter.LabelStatus.Bool)
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
