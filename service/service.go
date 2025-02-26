package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func InitServices() {
	ExampleService = NewExample()
}

// MAP TENANT_ID SQL_CONN
var MapDBConn map[string]sqlclient.ISqlClientConn

type DBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func NewDBConn(tenantId string, config DBConfig) (dbConn sqlclient.ISqlClientConn, err error) {
	sqlClientConfig := sqlclient.SqlConfig{
		Host:         config.Host,
		Port:         config.Port,
		Database:     config.Database,
		Username:     config.Username,
		Password:     config.Password,
		DialTimeout:  20,
		ReadTimeout:  30,
		WriteTimeout: 30,
		Timeout:      30,
		PoolSize:     10,
		MaxIdleConns: 10,
		MaxOpenConns: 10,
		Driver:       sqlclient.POSTGRESQL,
	}
	dbConn = sqlclient.NewSqlClient(sqlClientConfig)
	if err = dbConn.Connect(); err != nil {
		log.Error(err)
		err = ERR_DB_CONN_FAIL
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		repository.InitTables(ctx, dbConn)
	}()
	wg.Wait()
	MapDBConn[tenantId] = dbConn
	return
}

var ERR_EMPTY_CONN = errors.New("empty_conn")
var ERR_DB_CONN_FAIL = errors.New("db_conn_fail")

func GetDBConnOfUser(user model.AuthUser) (dbConn sqlclient.ISqlClientConn, err error) {
	if len(user.DatabaseHost) < 1 {
		err = ERR_EMPTY_CONN
		return
	}
	dbConn, ok := MapDBConn[user.TenantId]
	if !ok {
		dbConn, err = NewDBConn(user.TenantId, DBConfig{
			Host:     user.DatabaseHost,
			Port:     user.DatabasePort,
			Database: user.DatabaseName,
			Username: user.DatabaseUser,
			Password: user.DatabasePassword,
		})
		return
	}
	return
}
