package repository

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

var DBConn sqlclient.ISqlClientConn
var ESClient elasticsearch.IElasticsearchClient

func CreateTable(ctx context.Context, db sqlclient.ISqlClientConn, entity any) (err error) {
	_, err = db.GetDB().NewCreateTable().Model(entity).
		IfNotExists().
		Exec(ctx)
	return
}

func InitRepositories() {
	ExampleRepo = NewExample()
	AuthSourceRepo = NewAuthSource()
	ChatAppRepo = NewChatApp()
	ChatConnectionAppRepo = NewConnectionApp()
	ChatQueueRepo = NewChatQueue()
	ChatQueueAgentRepo = NewChatQueueAgent()
	ChatRoutingRepo = NewChatRouting()
	AgentAllocationRepo = NewAgentAllocation()
	ConnectionQueueRepo = NewConnectionQueue()
	ShareInfoRepo = NewShareInfo()
}

func InitRepositoriesES() {
	ESRepo = NewES()
	ConversationESRepo = NewConversationES()
	MessageESRepo = NewMessageES()
}

func InitTables(ctx context.Context, dbConn sqlclient.ISqlClientConn) {
	if err := CreateTable(ctx, dbConn, (*model.Example)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.AuthSource)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatApp)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatConnectionApp)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatQueue)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatRouting)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatQueueAgent)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.AgentAllocation)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ConnectionQueue)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ShareInfoForm)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.FacebookPage)(nil)); err != nil {
		log.Error(err)
	}
	log.Println("TABLES WERE CREATED")
}

func InitColumn(ctx context.Context, db sqlclient.ISqlClientConn) {
	if _, err := db.GetDB().NewAddColumn().Model((*model.ChatApp)(nil)).IfNotExists().ColumnExpr("default_app text null").Exec(ctx); err != nil {
		log.Info(err)
		panic(err)
	}
	if _, err := db.GetDB().NewAddColumn().Model((*model.AgentAllocation)(nil)).IfNotExists().ColumnExpr("app_id text").Exec(ctx); err != nil {
		log.Info(err)
		panic(err)
	}
	if _, err := db.GetDB().NewAddColumn().Model((*model.AgentAllocation)(nil)).IfNotExists().ColumnExpr("main_allocate text not null default active").Exec(ctx); err != nil {
		log.Info(err)
		panic(err)
	}
	if _, err := db.GetDB().NewAddColumn().Model((*model.AgentAllocation)(nil)).IfNotExists().ColumnExpr("source text").Exec(ctx); err != nil {
		log.Info(err)
		panic(err)
	}
}
