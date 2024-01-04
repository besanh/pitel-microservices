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
		// GetClient() *asynq.Client
	}

	QueueTask struct {
		RedisUrl  string
		Uri       string
		QueueName string
		MaxRetry  int
		Timeout   time.Duration
	}

	QueueTaskClient struct {
		Config QueueTask
		Client *asynq.Client
	}
)

const (
	MESSAGE_OTT_DELIVERY string = "message:delivery"
)

var QueueConnector *IQueueTaskClient

func NewQueueTaskClient(config QueueTask) IQueueTaskClient {
	queueTask := &QueueTaskClient{
		Config: config,
	}

	return queueTask
}

func (s *QueueTask) NewMessageOttDeliveryTask(agentId string, message *model.OttMessage) (*asynq.Task, error) {
	payload, err := json.Marshal(MessageDeliveryPayload{
		OttMessage: message,
		AgentId:    agentId,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(MESSAGE_OTT_DELIVERY, payload, asynq.MaxRetry(s.MaxRetry), asynq.Timeout(s.Timeout)), nil
}

func (s *QueueTask) HandleMessageDeliveryTask(ctx context.Context, task *asynq.Task) (*MessageDeliveryPayload, error) {
	var payload MessageDeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return &payload, nil
}
