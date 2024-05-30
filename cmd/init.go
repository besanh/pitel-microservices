package cmd

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/queue"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/messagequeue"
	"github.com/tel4vn/fins-microservices/internal/queuetask"
	"github.com/tel4vn/fins-microservices/internal/rabbitmq"
	"github.com/tel4vn/fins-microservices/internal/redis"
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

func initConfig() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := Config{
		Port:     env.GetStringENV("PORT", "8000"),
		LogLevel: env.GetStringENV("LOG_LEVEL", "error"),
		LogFile:  env.GetStringENV("LOG_FILE", "log/console.log"),
	}

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
