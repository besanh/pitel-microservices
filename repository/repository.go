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

func InitRepositories(db sqlclient.ISqlClientConn) {
	ExampleRepo = NewExample()
	IBKUserRepo = NewIBKUser(db)
	IBKTenantRepo = NewIBKTenant(db)
	IBKRoleRepo = NewIBKRole(db)
	IBKBusinessUnitRepo = NewIBKBusinessUnit(db)
}

func InitEsRepositories() {
	ESRepo = NewES()
}

func InitTables(ctx context.Context, dbConn sqlclient.ISqlClientConn) {
	if err := CreateTable(ctx, dbConn, (*model.IBKUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.IBKTenant)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.IBKRole)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.IBKBusinessUnit)(nil)); err != nil {
		log.Error(err)
	}
}
