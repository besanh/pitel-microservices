package repository

import (
	"context"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"time"
)

type (
	IChatAutoScriptToChatScript interface {
		Insert(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatAutoScriptToChatScript) error
		Update(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatAutoScriptToChatScript) error
	}

	ChatAutoScriptToChatScript struct{}
)

var ChatAutoScriptToChatScriptRepo IChatAutoScriptToChatScript

func NewChatAutoScriptToChatScript() IChatAutoScriptToChatScript {
	return &ChatAutoScriptToChatScript{}
}

func (c *ChatAutoScriptToChatScript) Insert(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatAutoScriptToChatScript) (err error) {
	entity.CreatedAt = time.Now()
	_, err = db.GetDB().NewInsert().
		Model(&entity).
		Exec(ctx)
	return err
}

func (c *ChatAutoScriptToChatScript) Update(ctx context.Context, db sqlclient.ISqlClientConn, entity model.ChatAutoScriptToChatScript) (err error) {
	entity.UpdatedAt = time.Now()
	_, err = db.GetDB().NewUpdate().
		Model(&entity).
		Where("chat_auto_script_id = ? and chat_script_id = ?", entity.ChatAutoScriptId, entity.ChatScriptId).
		Exec(ctx)
	return err
}
