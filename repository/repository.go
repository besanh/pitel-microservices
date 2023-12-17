package repository

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

var DBConn sqlclient.ISqlClientConn
var ES elasticsearch.IESClient
var ESClient elasticsearch.IElasticsearchClient

func CreateTable(ctx context.Context, db sqlclient.ISqlClientConn, entity any) (err error) {
	_, err = db.GetDB().NewCreateTable().Model(entity).
		IfNotExists().
		Exec(ctx)
	return
}

func InitRepositories() {
	ExampleRepo = NewExample()
	PluginConfigRepo = NewPluginConfig()
	RecipientConfigRepo = NewRecipientConfig()
	BalanceConfigRepo = NewBalanceConfig()
	TemplateBssRepo = NewTemplateBss()
	RoutingConfigRepo = NewRoutingConfig()
}

func InitEsRepositories() {
	InboxMarketingESRepo = NewInboxMarketingES()
	ESRepo = NewES()
}

func InitTables(ctx context.Context, dbConn sqlclient.ISqlClientConn) {
	if err := CreateTable(ctx, dbConn, (*model.Example)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.PluginConfig)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.RecipientConfig)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.BalanceConfig)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.TemplateBss)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.RoutingConfig)(nil)); err != nil {
		log.Error(err)
	}
}
