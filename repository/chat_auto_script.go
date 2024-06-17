package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/uptrace/bun"
	"time"
)

type (
	IChatAutoScript interface {
		IRepo[model.ChatAutoScript]
		InsertChatAutoScript(ctx context.Context, db sqlclient.ISqlClientConn, chatAutoScript model.ChatAutoScript, scripts []model.ChatAutoScriptToChatScript, labels []model.ChatAutoScriptToChatLabel) error
		GetChatAutoScripts(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAutoScriptFilter, limit, offset int) (int, *[]model.ChatAutoScriptView, error)
		GetChatAutoScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatAutoScriptView, error)
		DeleteChatAutoScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) error
	}

	ChatAutoScript struct {
		Repo[model.ChatAutoScript]
	}
)

var ChatAutoScriptRepo IChatAutoScript

func NewChatAutoScript() IChatAutoScript {
	return &ChatAutoScript{}
}

func (repo *ChatAutoScript) InsertChatAutoScript(ctx context.Context, db sqlclient.ISqlClientConn, chatAutoScript model.ChatAutoScript, scripts []model.ChatAutoScriptToChatScript, labels []model.ChatAutoScriptToChatLabel) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Error(errors.New("tx rollback failed"))
			}
		}
	}()

	chatAutoScript.CreatedAt = time.Now()
	_, err = tx.NewInsert().Model(&chatAutoScript).Exec(ctx)
	if err != nil {
		return err
	}

	for _, script := range scripts {
		if _, err = tx.NewInsert().Model(&script).Exec(ctx); err != nil {
			return err
		}
	}

	for _, label := range labels {
		if _, err = tx.NewInsert().Model(&label).Exec(ctx); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (repo *ChatAutoScript) GetChatAutoScripts(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatAutoScriptFilter, limit, offset int) (int, *[]model.ChatAutoScriptView, error) {
	result := new([]model.ChatAutoScriptView)
	query := db.GetDB().NewSelect().Model(result).
		Column("cas.*").
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name", "oa_info")
		}).
		Relation("ChatScriptLink", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("cas_cst.order ASC").
				Limit(3)
		}).
		Relation("ChatScriptLink.ChatScript", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("id")
		}).
		Relation("ChatLabelLink", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("cas_cl.order ASC").
				Limit(3)
		}).
		Relation("ChatLabelLink.ChatLabel", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("id")
		})
	if len(filter.ScriptName) > 0 {
		query.Where("cas.script_name ILIKE ?", "%"+filter.ScriptName+"%")
	}
	if len(filter.Channel) > 0 {
		query.Where("cas.channel = ?", filter.Channel)
	}
	if len(filter.OaId) > 0 {
		query.Where("connection_app.oa_info->cas.channel::text->0->>'oa_id' = ?", filter.OaId)
	}

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

func (repo *ChatAutoScript) GetChatAutoScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.ChatAutoScriptView, error) {
	result := new(model.ChatAutoScriptView)
	err := db.GetDB().NewSelect().
		Model(result).
		Relation("ConnectionApp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("connection_name")
		}).
		Relation("ChatScriptLink", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("cas_cst.order ASC").
				Limit(3)
		}).
		Relation("ChatScriptLink.ChatScript", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("id")
		}).
		Relation("ChatLabelLink", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("cas_cl.order ASC").
				Limit(3)
		}).
		Relation("ChatLabelLink.ChatLabel", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("id")
		}).
		Where("cas.id = ?", id).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *ChatAutoScript) DeleteChatAutoScriptById(ctx context.Context, db sqlclient.ISqlClientConn, id string) error {
	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Error(errors.New("tx rollback failed"))
			}
		}
	}()

	// delete related rows
	_, err = tx.NewDelete().Model((*model.ChatAutoScriptToChatScript)(nil)).
		Where("chat_auto_script_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewDelete().Model((*model.ChatAutoScriptToChatLabel)(nil)).
		Where("chat_auto_script_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	// delete the row
	_, err = tx.NewDelete().Model((*model.ChatAutoScript)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
