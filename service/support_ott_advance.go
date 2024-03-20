package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GetManageQueueAgent(ctx context.Context, queueId string) (manageQueueAgent *model.ChatManageQueueAgent, err error) {
	manageQueueAgentCache := cache.RCache.Get(MANAGE_QUEUE_AGENT + queueId)
	if manageQueueAgentCache != nil {
		if err = json.Unmarshal([]byte(manageQueueAgentCache.(string)), &manageQueueAgent); err != nil {
			log.Error(err)
		}
	}
	manageQueueAgent, err = repository.ManageQueueRepo.GetById(ctx, repository.DBConn, queueId)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
