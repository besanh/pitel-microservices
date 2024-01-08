package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func CheckChatQueueSetting(ctx context.Context, filter model.QueueFilter) (string, error) {
	var agentId string
	authUser := auth.ParseHeaderToUserDev(ctx)
	db, err := GetDBConnOfUser(*authUser)
	if err != nil {
		log.Error(err)
		return agentId, err
	}

	queue := model.ChatQueue{}
	// Get chat queue from cache or db
	queuesCache, err := cache.RCache.HGet(CHAT_QUEUE, filter.AppId)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else if len(queuesCache) < 1 {
		total, queues, err := repository.ChatQueueRepo.GetQueue(ctx, db, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if total < 1 {
			log.Error("queue not found")
			return agentId, errors.New("queue not found")
		}
		queue = (*queues)[0]
		if err := cache.RCache.RPush(ctx, CHAT_QUEUE, queue); err != nil {
			log.Error(err)
			return agentId, err
		}
	} else {
		err := json.Unmarshal([]byte(queuesCache), &queue)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
	}

	routing := model.ChatRouting{}
	// Get routing from cache or db
	routingCache, err := cache.RCache.HGet(CHAT_ROUTING, queue.Id)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else if len(routingCache) < 1 {
		routing, err := repository.ChatRoutingRepo.GetById(ctx, db, queue.ChatRoutingId)
		if err != nil {
			log.Error(err)
			return agentId, err
		} else if routing == nil {
			log.Error("routing not found")
			return agentId, errors.New("routing not found")
		}
		if err := cache.RCache.RPush(ctx, CHAT_ROUTING, routing); err != nil {
			log.Error(err)
			return agentId, err
		}
	} else {
		err := json.Unmarshal([]byte(routingCache), &routing)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
	}

	agents := []model.ChatQueueAgent{}
	// Chat queue agent
	filterChatQueueAgent := model.ChatQueueAgentFilter{
		QueueId: queue.Id,
	}
	chatQueueAgentCache, err := cache.RCache.HGet(CHAT_QUEUE_AGENT, queue.Id)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else if len(chatQueueAgentCache) < 1 {
		total, agents, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, db, filterChatQueueAgent, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if total < 1 {
			log.Error("agent not found")
			return agentId, errors.New("agent not found")
		}
		if err := cache.RCache.RPush(ctx, CHAT_QUEUE_AGENT, agents); err != nil {
			log.Error(err)
			return agentId, err
		}
	} else {
		err := json.Unmarshal([]byte(chatQueueAgentCache), &agents)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
	}

	if routing.RoutingName == "random" {
		rand.NewSource(time.Now().UnixNano())
		randomIndex := rand.Intn(len(agents))
		agent := agents[randomIndex]
		agentId = agent.AgentId
	} else if routing.RoutingName == "min_conversation" {
	}

	return agentId, nil
}

func GetConversationExist(ctx context.Context, data model.OttMessage) (conversation model.Conversation, err error) {
	conversationCache, err := cache.RCache.HGet(CONVERSATION, data.UserId)
	if err != nil {
		log.Error(err)
		return conversation, err
	} else if len(conversationCache) < 1 {
		conversation := model.Conversation{
			UserIdByApp: data.UserIdByApp,
			Username:    data.Username,
			Avatar:      data.Avatar,
		}
		err := cache.RCache.SADD(ctx, CONVERSATION, conversation)
		if err != nil {
			log.Error(err)
			return conversation, err
		}
		// Insert es
		_, err = InsertConversation(ctx, conversation)
		if err != nil {
			log.Error(err)
			return conversation, err
		}
	} else {
		err := json.Unmarshal([]byte(conversationCache), &conversation)
		if err != nil {
			log.Error(err)
			return conversation, err
		}
	}

	return conversation, nil
}

func InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error) {
	id = uuid.NewString()
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return id, err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return id, err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX, conversation.AppId); err != nil {
		log.Error(err)
		return id, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX, conversation.AppId); err != nil {
			log.Error(err)
			return id, err
		}
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.AppId, ES_INDEX, id, esDoc); err != nil {
		log.Error(err)
		return id, err
	}
	return id, nil
}

func PublishMessage(id string, message any) error {
	WsSubscribers.SubscribersMu.Lock()
	defer WsSubscribers.SubscribersMu.Unlock()
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	WsSubscribers.PublishLimiter.Wait(context.Background())
	isExisted := false
	for s := range WsSubscribers.Subscribers {
		if s.Id == id {
			isExisted = true
			select {
			case s.Message <- msgBytes:
			default:
				go s.CloseSlow()
			}
		}
	}
	if !isExisted {
		return errors.New("subscriber is not existed")
	}
	return nil
}
