package server

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/queuetask"
)

func NewMuxServer() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": env.GetIntENV("QUEUE_TASK_CRITICAL", 10),
				"default":  env.GetIntENV("QUEUE_TASK_DEFAULT", 5),
				"low":      env.GetIntENV("QUEUE_TASK_LOW", 1),
			},
		},
	)
	configQueueTask := queuetask.QueueTask{
		RedisUrl: env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		MaxRetry: env.GetIntENV("QUEUE_TASK_MAX_RETRY", 3),
		Timeout:  env.GetTimeDurationENV("QUEUE_TASK_TIMEOUT", 30*time.Second),
	}
	NewServeMuxQueueTask(srv, configQueueTask)
}

// Worker process queue task background
func NewServeMuxQueueTask(srv *asynq.Server, config queuetask.QueueTask) {
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()

	// taskMessage, err := queuetask.NewQueueTaskClient(config).NewMessageOttDeliveryTask()
	// if
	// taskHandleMessage, err := queuetask.NewQueueTaskClient(config).HandleMessageDeliveryTask(ctx, )

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	// mux.HandleFunc(queuetask.MESSAGE_OTT_DELIVERY)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
