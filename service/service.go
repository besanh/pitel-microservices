package service

import (
	"sync"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
)

func InitServices() {
	AuthService = NewAuthService("")
}

type DBConnection struct {
	sync.RWMutex
	MapConnection map[string]sqlclient.ISqlClientConn
}

func (db *DBConnection) GetConnection(tenant string) (client sqlclient.ISqlClientConn, err error) {
	db.RLock()
	defer db.RUnlock()
	client = db.MapConnection[tenant]
	if client == nil {
		err = variables.NewError(constants.ERR_DB_CONNECTION_ERROR)
		return
	}
	return
}

func (db *DBConnection) SetConnection(tenant string, client sqlclient.ISqlClientConn) {
	db.Lock()
	defer db.Unlock()
	db.MapConnection[tenant] = client
}

var DBCons DBConnection

func InitDBConnection(conn sqlclient.ISqlClientConn) {
	DBCons = DBConnection{
		MapConnection: make(map[string]sqlclient.ISqlClientConn),
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// defer cancel()

	// // Get tenants
	// tenants, _, err := repository.Tem.SelectByQueryWithDB(ctx, conn, nil, 99, 0, "")
	// if err != nil {
	// 	log.Errorf("[MigrateDBCollections] Error get tenants: %v", err)
	// 	return
	// }

	// for _, tenant := range tenants {
	// 	log.Infof("[InitDBConnection] Get database setting for tenant: %v", tenant.TenantName)
	// 	var dbCfg model.AAA_DB
	// 	if err = util.ParseAnyToAny(tenant.DatabaseSetting, &dbCfg); err != nil {
	// 		log.Errorf("[InitDBConnection] Error parse database setting: %v", err)
	// 		continue
	// 	}
	// 	log.Infof("[InitDBConnection] Get collections for tenant: %v - database: %v", tenant.TenantName, dbCfg.DBName)
	// 	// Handle check database exist
	// 	var dbConfig sqlclient.SqlConfig
	// 	var db sqlclient.ISqlClientConn
	// 	dbConfig = sqlclient.SqlConfig{
	// 		Host:         dbCfg.Host,
	// 		Database:     dbCfg.DBName,
	// 		Username:     dbCfg.Username,
	// 		Password:     dbCfg.Password,
	// 		Port:         dbCfg.Port,
	// 		DialTimeout:  20,
	// 		ReadTimeout:  30,
	// 		WriteTimeout: 30,
	// 		Timeout:      30,
	// 		PoolSize:     10,
	// 		MaxOpenConns: 20,
	// 		MaxIdleConns: 10,
	// 		Driver:       sqlclient.POSTGRESQL,
	// 	}
	// 	db, err = sqlclient.NewSqlClient(dbConfig)
	// 	if err != nil {
	// 		log.Errorf("[InitDBConnection] Error connect database: %v", err)
	// 		continue
	// 	}
	// 	DBCons.SetConnection(tenant.Id, db)
	// }
}
