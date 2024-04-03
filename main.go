package main

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/queue"
	elasticsearch "github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/messagequeue"
	"github.com/tel4vn/fins-microservices/internal/queuetask"
	"github.com/tel4vn/fins-microservices/internal/rabbitmq"
	"github.com/tel4vn/fins-microservices/internal/redis"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/internal/storage"
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

var (
	config Config
)

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
	// rabbitmqconfig := rmq.Config{
	// 	Uri:                  env.GetStringENV("RMQ_HOST", "rabbitmq.dev.fins.vn"),
	// 	ChannelNotifyTimeout: 100 * time.Millisecond,
	// 	Reconnect: struct {
	// 		Interval   time.Duration
	// 		MaxAttempt int
	// 	}{
	// 		Interval:   500 * time.Millisecond,
	// 		MaxAttempt: 7200,
	// 	},
	// }
	// rmq.RabbitConnector = rmq.New(rabbitmqconfig)
	// rmq.RabbitConnector.RoutingKey = "es.writer"
	// rmq.RabbitConnector.ExchangeName = "events"
	// if err := rmq.RabbitConnector.Ping(); err != nil {
	// 	panic(err)
	// }
	rabbitMQConfig := rabbitmq.Config{
		Addr:         env.GetStringENV("RMQ_HOST", "rabbitmq.dev.fins.vn"),
		ExchangeName: env.GetStringENV("RMQ_EXCHANGE_NAME", "bss-message"),
		QueueName:    env.GetStringENV("RMQ_QUEUE_NAME", "bss-chat"),
	}

	err = messagequeue.NewMQConn(rabbitMQConfig)
	if err != nil {
		panic(err)
	}
	esCfg := elasticsearch.Config{
		Username:              env.GetStringENV("ES_USERNAME", "elastic"),
		Password:              env.GetStringENV("ES_PASSWORD", "tel4vnEs2021"),
		Host:                  []string{env.GetStringENV("ES_HOST", "http://113.164.246.12:9200")},
		MaxRetries:            10,
		ResponseHeaderTimeout: 60,
		RetryStatuses:         []int{502, 503, 504},
	}
	repository.ESClient = elasticsearch.NewElasticsearchClient(esCfg)

	// Queue task
	queueTaskConfig := queuetask.QueueTask{
		RedisUrl: env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		MaxRetry: env.GetIntENV("QUEUE_TASK_MAX_RETRY", 3),
		Timeout:  env.GetTimeDurationENV("QUEUE_TASK_TIMEOUT", 30*time.Second),
	}
	queuetask.NewQueueTaskClient(queueTaskConfig)

	// goauth.GoAuthClient = goauth.NewGoAuth(cache.RCache.GetClient())
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
	cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())
	defer cache.RCache.Close()

	cache.MCache = cache.NewMemCache()
	defer cache.MCache.Close()

	// Init Repositories
	repository.InitRepositories()
	repository.InitRepositoriesES()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repository.InitTables(ctx, repository.DBConn)
	repository.InitColumn(ctx, repository.DBConn)

	// Init services
	service.MapDBConn = make(map[string]sqlclient.ISqlClientConn, 0)
	service.ES_INDEX = env.GetStringENV("ES_INDEX", "pitel_bss_chat")
	service.ES_INDEX_CONVERSATION = env.GetStringENV("ES_INDEX_CONVERSATION", "pitel_bss_conversation")
	service.OTT_URL = env.GetStringENV("OTT_DOMAIN", "")
	service.OTT_VERSION = env.GetStringENV("OTT_VERSION", "v1")
	service.API_SHARE_INFO_HOST = env.GetStringENV("API_SHARE_INFO_HOST", "https://api.dev.fins.vn")
	service.API_CRM = env.GetStringENV("API_CRM", "")
	service.ENABLE_PUBLISH_ADMIN = env.GetBoolENV("ENABLE_PUBLISH_ADMIN", false)
	service.AAA_HOST = env.GetStringENV("AAA_HOST", "https://aaa.dev.fins.vn")
	service.InitServices()

	// Init storage
	storage.InitStorage()

	// Run cron jobs
	// handleCronBatchSchedule(service.BatchService.ScanBatchJobEvery1Minute)

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

func handleCronBatchSchedule(f func()) {
	s1 := gocron.NewScheduler(time.UTC)
	s1.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err := s1.Every(1).Minute().Do(f)
	if err != nil {
		return
	}
	s1.StartAsync()
}
