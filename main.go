package main

import (
	"io"
	"path/filepath"
	"time"

	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/queue"
	elasticsearchsearch "github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/queuetask"
	"github.com/tel4vn/fins-microservices/internal/redis"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/server"
	"github.com/tel4vn/fins-microservices/service"

	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port       string
	gRPCPort   string
	AAA_Adress string
	LogLevel   string
	LogFile    string
}

var config Config

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := Config{
		Port:       env.GetStringENV("PORT", "8000"),
		gRPCPort:   env.GetStringENV("GRPC_PORT", "8001"),
		AAA_Adress: env.GetStringENV("AAA_ADRESS", "aaa-service:8001"),
		LogLevel:   env.GetStringENV("LOG_LEVEL", "error"),
		LogFile:    env.GetStringENV("LOG_FILE", "log/console.log"),
	}

	sqlClientConfig := sqlclient.SqlConfig{
		Host:         env.GetStringENV("DB_HOST", "localhost"),
		Database:     env.GetStringENV("DB_DATABASE", "dev_fins_aaa"),
		Username:     env.GetStringENV("DB_USERNAME", "admin"),
		Password:     env.GetStringENV("DB_PASSWORD", "password"),
		Port:         env.GetIntENV("DB_PORT", 5432),
		DialTimeout:  20,
		ReadTimeout:  30,
		WriteTimeout: 30,
		Timeout:      30,
		PoolSize:     10,
		MaxOpenConns: 20,
		MaxIdleConns: 10,
		Driver:       sqlclient.POSTGRESQL,
	}
	repository.DBConn = sqlclient.NewSqlClient(sqlClientConfig)
	var err error
	if redis.Redis, err = redis.NewRedis(redis.Config{
		Addr:         env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		Password:     env.GetStringENV("REDIS_PASSWORD", ""),
		DB:           env.GetIntENV("REDIS_DATABASE", 0),
		PoolSize:     30,
		PoolTimeout:  20,
		IdleTimeout:  10,
		ReadTimeout:  20,
		WriteTimeout: 15,
	}); err != nil {
		panic(err)
	}
	queue.RMQ = queue.NewRMQ(queue.Rcfg{
		Address:  env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		Password: env.GetStringENV("REDIS_PASSWORD", ""),
		DB:       9,
	})
	esCfg := elasticsearchsearch.Config{
		Username:              env.GetStringENV("ES_USERNAME", "elastic"),
		Password:              env.GetStringENV("ES_PASSWORD", "tel4vnEs2021"),
		Host:                  []string{env.GetStringENV("ES_HOST", "http://113.164.246.12:9200")},
		MaxRetries:            10,
		ResponseHeaderTimeout: 60,
		RetryStatuses:         []int{502, 503, 504},
	}
	repository.ESClient = elasticsearchsearch.NewElasticsearchClient(esCfg)

	// Queue task
	queueTaskConfig := queuetask.QueueTask{
		RedisUrl: env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		MaxRetry: env.GetIntENV("QUEUE_TASK_MAX_RETRY", 3),
		Timeout:  env.GetTimeDurationENV("QUEUE_TASK_TIMEOUT", 30*time.Second),
	}
	queuetask.NewQueueTaskClient(queueTaskConfig)

	// goauth.GoAuthClient = goauth.NewGoAuth(redis.Redis.GetClient())
	// authMdw.SetupGoGuardian()
	authMdw.AuthMdw = authMdw.NewGatewayAuthMiddleware(env.GetStringENV("ENV", "dev"))
	config = cfg
}

func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)

	// Init gRPC client
	// auth.GRPC_Client = auth.NewGRPCAuh(config.AAA_Adress)

	// Init Repositories
	repository.InitRepositories()
	repository.InitRepositoriesES()

	// Init services
	service.MapDBConn = make(map[string]sqlclient.ISqlClientConn, 0)
	service.InitServices()

	// Run gRPC server
	server.NewGRPCServer(config.gRPCPort)
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
}
