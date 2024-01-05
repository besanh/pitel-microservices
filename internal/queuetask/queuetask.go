package queuetask

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IQueueTaskClient interface {
		GetClient() *asynq.Client
		NewMessageOttDeliveryTask(message *model.OttMessage) (*asynq.Task, error)
		HandleMessageDeliveryTask(ctx context.Context, task *asynq.Task) (*MessageDeliveryPayload, error)
	}

	QueueTask struct {
		RedisUrl  string
		Uri       string
		QueueName string
		MaxRetry  int
		Timeout   time.Duration
	}

	QueueTaskClient struct {
		config QueueTask
		client *asynq.Client
	}
)

const (
	MESSAGE_OTT_DELIVERY string = "message:delivery"
)

var QueueConnector *IQueueTaskClient

func NewQueueTaskClient(config QueueTask) IQueueTaskClient {
	queueTask := &QueueTaskClient{
		config: config,
	}

	return queueTask
}

func (s *QueueTaskClient) GetClient() *asynq.Client {
	return s.client
}

func (s *QueueTaskClient) NewMessageOttDeliveryTask(message *model.OttMessage) (*asynq.Task, error) {
	payload, err := json.Marshal(MessageDeliveryPayload{
		OttMessage: message,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		MESSAGE_OTT_DELIVERY,
		payload,
		asynq.MaxRetry(s.config.MaxRetry),
		asynq.Timeout(s.config.Timeout),
	), nil
}

func (s *QueueTaskClient) HandleMessageDeliveryTask(ctx context.Context, task *asynq.Task) (*MessageDeliveryPayload, error) {
	var payload MessageDeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return &payload, nil
}

func (s *QueueTaskClient) EnQueue(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	info, err := s.client.Enqueue(task, asynq.Queue(s.config.QueueName))
	if err != nil {
		return nil, err
	}
	return info, nil
}
