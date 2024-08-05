package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
)

type (
	IConversationMedia interface {
		IRepo[model.ConversationMedia]
		GetConversationMedias(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConversationMediaFilter, limit, offset int) (int, *[]model.ConversationMedia, error)
	}
	ConversationMedia struct {
		Repo[model.ConversationMedia]
	}
)

var ConversationMediaRepo IConversationMedia

func NewConversationMedia() IConversationMedia {
	return &ConversationMedia{}
}

func (repo *ConversationMedia) GetConversationMedias(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConversationMediaFilter, limit, offset int) (int, *[]model.ConversationMedia, error) {
	result := new([]model.ConversationMedia)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query = query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConversationId) > 0 {
		query = query.Where("conversation_id = ?", filter.ConversationId)
	}
	if len(filter.ExternalConversationId) > 0 {
		query = query.Where("external_conversation_id = ?", filter.ExternalConversationId)
	}
	if len(filter.ConversationType) > 0 {
		query = query.Where("conversation_type = ?", filter.ConversationType)
	}
	if len(filter.MediaType) > 0 {
		query = query.Where("media_type = ?", filter.MediaType)
	}
	if len(filter.MediaName) > 0 {
		query.Where(
			"? ILIKE ? OR ? ILIKE ?",
			bun.Ident("media_header"), "%"+filter.MediaName+"%",
			bun.Ident("media_url"), "%"+filter.MediaName+"%",
		)
	}
	query.Order("send_timestamp DESC", "created_at DESC")

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
