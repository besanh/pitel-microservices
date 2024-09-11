package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
)

type (
	INotesList interface {
		IRepo[model.NotesList]
		GetNotesLists(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConversationNotesListFilter, limit, offset int) (int, *[]model.NotesList, error)
	}
	NotesList struct {
		Repo[model.NotesList]
	}
)

var NotesListRepo INotesList

func NewNotesList() INotesList {
	return &NotesList{}
}

func (repo *NotesList) GetNotesLists(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ConversationNotesListFilter, limit, offset int) (int, *[]model.NotesList, error) {
	result := new([]model.NotesList)
	query := db.GetDB().NewSelect().Model(result)
	if len(filter.TenantId) > 0 {
		query.Where("tenant_id = ?", filter.TenantId)
	}
	if len(filter.ConversationId) > 0 {
		query.Where("conversation_id = ?", filter.ConversationId)
	}
	if len(filter.AppId) > 0 {
		query.Where("app_id = ?", filter.AppId)
	}
	if len(filter.OaId) > 0 {
		query.Where("oa_id = ?", filter.OaId)
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
