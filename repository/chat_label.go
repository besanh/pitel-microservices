package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
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
	if len(filter.AppId) > 0 {
		query.Where("app_id = ?", filter.AppId)
	}
	if len(filter.OaId) > 0 {
		query.Where("oa_id = ?", filter.OaId)
	}
	if len(filter.LabelType) > 0 {
		query.Where("label_type = ?", filter.LabelType)
	}
	if len(filter.LabelName) > 0 {
		if filter.IsSearchExactly.Valid && filter.IsSearchExactly.Bool {
			query.Where("label_name = ?", filter.LabelName)
		} else {
			query.Where("? = ?", bun.Ident("label_name"), filter.LabelName)
		}
	}
	if len(filter.LabelColor) > 0 {
		query.Where("label_color = ?", filter.LabelColor)
	}
	if filter.LabelStatus.Valid {
		query.Where("label_status = ?", filter.LabelStatus.Bool)
	}
	if len(filter.ExternalLabelId) > 0 {
		query.Where("external_label_id = ?", filter.ExternalLabelId)
	}
	if len(filter.LabelIds) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			listUuid := make([]string, 0)
			externalId := make([]string, 0)
			for _, item := range filter.LabelIds {
				if IsValidUUID(item) {
					listUuid = append(listUuid, item)
				} else {
					externalId = append(externalId, item)
				}
			}
			if len(listUuid) > 0 {
				q.Where("id IN (?)", bun.In(listUuid))
			}
			if len(externalId) > 0 {
				q.Where("external_label_id IN (?)", bun.In(externalId))
			}
			return q
		})
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
