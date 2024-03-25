package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GetManageQueueAgent(ctx context.Context, queueId string) (manageQueueAgent model.ChatManageQueueAgent, err error) {
	manageQueueAgentCache := cache.RCache.Get(MANAGE_QUEUE_AGENT + "_" + queueId)
	if manageQueueAgentCache != nil {
		if err = json.Unmarshal([]byte(manageQueueAgentCache.(string)), &manageQueueAgent); err != nil {
			log.Error(err)
		}
	}
	filter := model.ChatManageQueueAgentFilter{
		QueueId: queueId,
	}
	total, manageQueueAgents, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		manageQueueAgent = (*manageQueueAgents)[0]
		if err = cache.RCache.Set(MANAGE_QUEUE_AGENT+"_"+queueId, manageQueueAgent, MANAGE_QUEUE_AGENT_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}
	return manageQueueAgent, nil
}

func GenerateConversationId(appId, conversationId string) (newConversationId string) {
	newConversationId = appId + "_" + conversationId
	return
}

func RoundRobinAgentOnline(ctx context.Context, conversationId string, queueAgents *[]model.ChatQueueAgent) (*Subscriber, error) {
	userLive := Subscriber{}
	userLives := []Subscriber{}
	subscribers, err := cache.RCache.HGetAll(BSS_SUBSCRIBERS)
	if err != nil {
		log.Error(err)
		return &userLive, err
	}
	for _, item := range subscribers {
		s := Subscriber{}
		if err := json.Unmarshal([]byte(item), &s); err != nil {
			log.Error(err)
			return &userLive, err
		}
		if (s.Level == "user" || s.Level == "agent") && CheckInLive(*queueAgents, s.Id) {
			userLives = append(userLives, s)
		}
	}
	if len(userLives) > 0 {
		index, userAllocate := GetAgentIsRoundRobin(userLives)
		userLive = *userAllocate
		userLive.IsAssignRoundRobin = true
		userPrevious := Subscriber{}
		if index < len(userLives) {
			userPrevious = userLives[(index+1)%len(userLives)]
		} else {
			userPrevious = userLives[0]
		}
		userPrevious.IsAssignRoundRobin = false

		// Update current
		jsonByteUserLive, err := json.Marshal(&userLive)
		if err != nil {
			log.Error(err)
			return &userLive, err
		}
		if err := cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, userLive.Id, string(jsonByteUserLive)); err != nil {
			log.Error(err)
			return &userLive, err
		}

		// Update previous
		if userPrevious.Id != userLive.Id {
			jsonByteUserLivePrevious, err := json.Marshal(&userPrevious)
			if err != nil {
				log.Error(err)
				return &userLive, err
			}
			if err := cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, userPrevious.Id, string(jsonByteUserLivePrevious)); err != nil {
				log.Error(err)
				return &userLive, err
			}
		}
		return &userLive, nil
	} else {
		return &userLive, errors.New("no user online")
	}
}

func GetAgentIsRoundRobin(userLives []Subscriber) (int, *Subscriber) {
	isOk := false
	index := 0
	userLive := Subscriber{}
	for i, item := range userLives {
		if item.IsAssignRoundRobin {
			if (i+1)%len(userLives) <= len(userLives) {
				userLive = userLives[(i+1)%len(userLives)]
				isOk = true
				index = (i + 1) % len(userLives)
				break
			} else {
				isOk = true
				userLive = userLives[0]
				break
			}
		}
	}
	if isOk {
		return index, &userLive
	}
	userLive = userLives[0]
	return index, &userLive
}

func CheckInLive(queueAgents []model.ChatQueueAgent, id string) bool {
	for _, item := range queueAgents {
		if item.AgentId == id {
			return true
		}
	}
	return false
}

// TODO: caching
func CacheConnection(ctx context.Context, connectionId string, conversation model.Conversation) (model.Conversation, error) {
	connectionExist, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, connectionId)
	if err != nil {
		log.Error(err)
		return conversation, err
	}
	if connectionExist != nil {
		if conversation.ConversationType == "zalo" {
			conversation.OaName = connectionExist.OaInfo.Zalo[0].OaName
			conversation.OaAvatar = connectionExist.OaInfo.Zalo[0].Avatar
		} else if conversation.ConversationType == "facebook" {
			conversation.OaName = connectionExist.OaInfo.Facebook[0].OaName
			conversation.OaAvatar = connectionExist.OaInfo.Facebook[0].Avatar
		}
	}
	return conversation, nil
}

// TODO: push event to admin
// func PushEventToAdmin(ctx context.Context)
