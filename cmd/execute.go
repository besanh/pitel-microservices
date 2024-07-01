package cmd

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/queue"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/queuetask"
	"github.com/tel4vn/fins-microservices/internal/rabbitmq"
	"github.com/tel4vn/fins-microservices/internal/redis"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/repository"
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
		Port:     env.GetStringENV("PORT", "8000"),
		gRPCPort: env.GetStringENV("GRPC_PORT", "8002"),
		LogLevel: env.GetStringENV("LOG_LEVEL", "error"),
		LogFile:  env.GetStringENV("LOG_FILE", "log/console.log"),
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

	// Init Redis Message Queue
	redisCfg := queue.Rcfg{
		Address:  env.GetStringENV("REDIS_ADDRESS", "localhost"),
		Password: env.GetStringENV("REDIS_PASSWORD", ""),
		DB:       env.GetIntENV("REDIS_RMQ_DATABASE", 9),
	}
	queue.RMQ = queue.NewRMQ(redisCfg)

	// RabbitMQ
	rabbitmqconfig := rabbitmq.Config{
		Uri:                  env.GetStringENV("RMQ_HOST", "amqp://guest:guest@localhost:5672/"),
		ChannelNotifyTimeout: 10 * time.Second,
		Reconnect: struct {
			Interval   time.Duration
			MaxAttempt int
		}{
			Interval:   500 * time.Millisecond,
			MaxAttempt: 7200,
		},
	}

	rabbitmq.RabbitConnector = rabbitmq.New(rabbitmqconfig)
	rabbitmq.RabbitConnector.RoutingKey = env.GetStringENV("RABBITMQ_ROUTING_KEY", "bss-chat")
	rabbitmq.RabbitConnector.ExchangeName = env.GetStringENV("RABBITMQ_EXCHANGE_NAME", "bss-message")
	if err := rabbitmq.RabbitConnector.Ping(); err != nil {
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

func Execute() {
	var rootCmd = cobra.Command{Use: "chat-service"}
	rootCmd.AddCommand(cmdMain)
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
