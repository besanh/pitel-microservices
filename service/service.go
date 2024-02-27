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
	"golang.org/x/time/rate"
)

func InitServices() {
	ExampleService = NewExample()
}

// MAP TENANT_ID SQL_CONN
var (
	MapDBConn        map[string]sqlclient.ISqlClientConn
	ERR_EMPTY_CONN   = errors.New("empty_conn")
	ERR_DB_CONN_FAIL = errors.New("db_conn_fail")

	// ES
	ES_HOST               = "https://es.dev.fins.vn"
	ES_USERNAME           = "elastic"
	ES_PASSWORD           = "FinS##TEL4VN##ES#!2324"
	ES_INDEX              = "" //             = "pitel_bss_chat"
	ES_INDEX_CONVERSATION = "" // = "pitel_bss_conversation"

	// Redis
	CONVERSATION            = "conversation"
	CONVERSATION_EXPIRE     = 30 * time.Minute
	CHAT_QUEUE              = "chat_queue"
	CHAT_QUEUE_EXPIRE       = 30 * time.Minute
	CHAT_ROUTING            = "chat_routing"
	CHAT_ROUTING_EXPIRE     = 1 * time.Hour
	CHAT_QUEUE_AGENT        = "chat_queue_agent"
	CHAT_QUEUE_AGENT_EXPIRE = 10 * time.Minute
	CHAT_APP                = "chat_app"
	CHAT_APP_EXPIRE         = 5 * time.Hour
	AGENT_ALLOCATION        = "agent_allocation"
	AGENT_ALLOCATION_EXPIRE = 1 * time.Hour

	ORIGIN_LIST = []string{"localhost:*", "*.tel4vn.com"}

	OTT_URL             string = ""
	API_SHARE_INFO_HOST string = ""
)

type (
	DBConfig struct {
		Host     string
		Port     int
		Database string
		Username string
		Password string
	}

	Subscriber struct {
		Id             string      `json:"id"`
		TenantId       string      `json:"tenant_id"`
		BusinessUnitId string      `json:"business_unit_id"`
		UserId         string      `json:"user_id"`
		Username       string      `json:"username"`
		Services       []string    `json:"services"`
		Level          string      `json:"level"`
		Message        chan []byte `json:"-"`
		CloseSlow      func()      `json:"-"`
	}

	Subscribers struct {
		SubscriberMessageBuffer int
		PublishLimiter          *rate.Limiter
		SubscribersMu           sync.Mutex
		Subscribers             map[*Subscriber]struct{}
	}
)

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

func HandleGetDBConSource(authUser *model.AuthUser) (sqlclient.ISqlClientConn, error) {
	var dbCon sqlclient.ISqlClientConn
	if len(authUser.Source) < 1 || authUser.Source == "authen" {
		dbCon = repository.DBConn
	} else {
		dbConTmp, err := GetDBConnOfUser(*authUser)
		if err != nil {
			return dbCon, err
		}
		dbCon = dbConTmp
	}
	return dbCon, nil
}
