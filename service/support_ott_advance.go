package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GetManageQueueAgent(ctx context.Context, queueId string) (manageQueueAgent model.ChatManageQueueAgent, err error) {
	manageQueueAgentCache := cache.RCache.Get(MANAGE_QUEUE_AGENT + queueId)
	if manageQueueAgentCache != nil {
		if err = json.Unmarshal([]byte(manageQueueAgentCache.(string)), &manageQueueAgent); err != nil {
			log.Error(err)
		}
	}
	filter := model.ManageQueueAgentFilter{
		QueueId: queueId,
	}
	total, manageQueueAgents, err := repository.ManageQueueRepo.GetManageQueue(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		manageQueueAgent = (*manageQueueAgents)[0]
		if err = cache.RCache.Set(MANAGE_QUEUE_AGENT+queueId, manageQueueAgent, MANAGE_QUEUE_AGENT_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}
	return manageQueueAgent, nil
}
