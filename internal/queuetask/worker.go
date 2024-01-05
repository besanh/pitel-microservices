package queuetask

import "github.com/hibiken/asynq"

type (
	IWorkerClient interface {
	}
	WorkerConfig struct {
		Url         string
		Concurrency int64
		Queues      map[string]int
	}
	WorkerClient struct {
		Config WorkerConfig
		Client *asynq.Client
	}
)

func NewWorkerClient(config WorkerConfig) IWorkerClient {
	worker := &WorkerClient{
		Config: config,
		Client: asynq.NewClient(asynq.RedisClientOpt{Addr: config.Url}),
	}

	return worker
}
