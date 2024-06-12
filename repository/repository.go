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
	ChatQueueUserRepo = NewChatQueueUser()
	ChatRoutingRepo = NewChatRouting()
	UserAllocateRepo = NewUserAllocate()
	ConnectionQueueRepo = NewConnectionQueue()
	ShareInfoRepo = NewShareInfo()
	ManageQueueRepo = NewManageQueue()
	ChatMsgSampleRepo = NewChatMsgSample()
	ChatPersonalizationRepo = NewChatPersonalization()
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
	if err := CreateTable(ctx, dbConn, (*model.ChatQueueUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.UserAllocate)(nil)); err != nil {
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
	if err := CreateTable(ctx, dbConn, (*model.ChatManageQueueUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatEmail)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatMsgSample)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatPersonalization)(nil)); err != nil {
		log.Error(err)
	}
	log.Println("TABLES WERE CREATED")
}

func InitColumn(ctx context.Context, db sqlclient.ISqlClientConn) {
	if _, err := db.GetDB().NewAddColumn().Model((*model.ChatApp)(nil)).IfNotExists().ColumnExpr("default_app text null").Exec(ctx); err != nil {
		log.Info(err)
		panic(err)
	}
	if _, err := db.GetDB().NewAddColumn().Model((*model.UserAllocate)(nil)).IfNotExists().ColumnExpr("oa_id text not null").Exec(ctx); err != nil {
		log.Info(err)
		panic(err)
	}
}

func InitRows(ctx context.Context, db sqlclient.ISqlClientConn) {
	total, err := db.GetDB().NewSelect().Model((*model.ChatPersonalization)(nil)).Count(ctx)
	if err != nil {
		log.Info(err)
		panic(err)
	}
	if total == 0 {
		if err := ChatPersonalizationRepo.InsertDefaultPersonalizationValue(ctx, db, "page_name"); err != nil {
			log.Info(err)
			panic(err)
		}
		if err := ChatPersonalizationRepo.InsertDefaultPersonalizationValue(ctx, db, "customer_name"); err != nil {
			log.Info(err)
			panic(err)
		}
		if err := ChatPersonalizationRepo.InsertDefaultPersonalizationValue(ctx, db, "gender"); err != nil {
			log.Info(err)
			panic(err)
		}
	}
}
