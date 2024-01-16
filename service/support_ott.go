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
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func CheckChatQueueSetting(ctx context.Context, filter model.QueueFilter, userIdByApp string) (string, error) {
	var agentId string
	chatQueue := model.ChatQueue{}

	chatQueueCache := cache.RCache.Get(CHAT_QUEUE + "_" + filter.AppId)
	if chatQueueCache != nil {
		if err := json.Unmarshal([]byte(chatQueueCache.(string)), &chatQueue); err != nil {
			log.Error(err)
			return agentId, err
		}
	} else {
		total, queues, err := repository.ChatQueueRepo.GetQueue(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if total < 1 {
			log.Error("queue not found")
			return agentId, errors.New("queue not found")
		}
		chatQueue = (*queues)[0]
		if err := cache.RCache.Set(CHAT_QUEUE+"_"+filter.AppId, chatQueue, CHAT_QUEUE_EXPIRE); err != nil {
			log.Error(err)
			return agentId, err
		}
	}

	routing := model.ChatRouting{}
	// Get routing from cache or db
	chatRoutingCache := cache.RCache.Get(CHAT_ROUTING + "_" + chatQueue.ChatRoutingId)
	if chatRoutingCache != nil {
		if err := json.Unmarshal([]byte(chatRoutingCache.(string)), &routing); err != nil {
			log.Error(err)
			return agentId, err
		}
	} else {
		routingTmp, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, chatQueue.ChatRoutingId)
		if err != nil {
			log.Error(err)
			return agentId, err
		} else if routingTmp == nil {
			log.Error("routing not found")
			return agentId, errors.New("routing not found")
		}
		if err := cache.RCache.Set(CHAT_ROUTING+"_"+chatQueue.ChatRoutingId, routingTmp, CHAT_ROUTING_EXPIRE); err != nil {
			log.Error(err)
			return agentId, err
		}
		routing = *routingTmp
	}

	if routing.RoutingName == "random" {
		subscribers := []Subscriber{}
		// Get subscriber random
		// Get subscriber having conversation
		// Assign conversation for subscriber
		rand.NewSource(time.Now().UnixNano())
		randomIndex := rand.Intn(len(WsSubscribers.Subscribers))
		if len(WsSubscribers.Subscribers) > 0 {
			for s := range WsSubscribers.Subscribers {
				subscribers = append(subscribers, *s)
			}
			agent := subscribers[randomIndex]
			filter := model.AgentAllocationFilter{
				UserIdByApp: userIdByApp,
			}
			total, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
			if err != nil {
				log.Error(err)
				return agentId, err
			}
			if total > 0 {
				for _, item := range *agentAllocations {
					if len(item.AgentId) > 1 {
						agentId = item.AgentId
						break
					}
				}
			}
			if len(agentId) < 1 {
				agentId = agent.UserId
				agentAllocation := model.AgentAllocation{
					Base:               model.InitBase(),
					UserIdByApp:        userIdByApp,
					AgentId:            agent.UserId,
					QueueId:            chatQueue.Id,
					AllocatedTimestamp: time.Now().Unix(),
				}
				if err := repository.AgentAllocationRepo.Insert(ctx, repository.DBConn, agentAllocation); err != nil {
					log.Error(err)
					return agentId, err
				}
			} else if len(agentId) > 0 {
				agentId = agent.UserId
				filter := model.AgentAllocationFilter{
					UserIdByApp: userIdByApp,
				}
				total, agentAllocationTmp, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
				if err != nil {
					log.Error(err)
					return agentId, err
				}
				if total > 0 {
					agentAllocation := (*agentAllocationTmp)[0]
					agentAllocation.AgentId = agent.UserId
					agentAllocation.QueueId = chatQueue.Id
					agentAllocation.AllocatedTimestamp = time.Now().Unix()
					if err := repository.AgentAllocationRepo.Update(ctx, repository.DBConn, agentAllocation); err != nil {
						log.Error(err)
						return agentId, err
					}
				}
			}
		}
	} else if routing.RoutingName == "min_conversation" {
	}

	return agentId, nil
}

func GetConversationExist(ctx context.Context, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
	conversation = model.Conversation{
		ConversationId:   uuid.NewString(),
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		UserIdByApp:      data.UserIdByApp,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		Uid:              data.UserId,
		CreatedAt:        time.Now().Format(time.RFC3339),
	}

	isExisted := false
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.UserIdByApp)
	if conversationCache != nil {
		isExisted = true
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		return conversation, isNew, nil
	} else {
		filter := model.ConversationFilter{
			AppId:       []string{data.AppId},
			UserIdByApp: []string{data.UserIdByApp},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, data.AppId, ES_INDEX_CONVERSATION, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if total > 0 {
			conversation = (*conversations)[0]
			if err := cache.RCache.Set(CONVERSATION+"_"+data.UserIdByApp, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}
	}

	if !isExisted {
		id, err := InsertConversation(ctx, conversation)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		conversation.ConversationId = id
		if err := cache.RCache.Set(CONVERSATION+"_"+data.UserIdByApp, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isNew = true
	}

	return conversation, isNew, nil
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
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX_CONVERSATION, conversation.AppId); err != nil {
		log.Error(err)
		return id, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX_CONVERSATION, conversation.AppId); err != nil {
			log.Error(err)
			return id, err
		}
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.AppId, ES_INDEX_CONVERSATION, id, esDoc); err != nil {
		log.Error(err)
		return id, err
	}
	return id, nil
}

func CheckConversationInAgent(userIdByApp string, allocationAgent []*model.AgentAllocation) bool {
	for _, item := range allocationAgent {
		if item.UserIdByApp == userIdByApp {
			return true
		}
	}
	return false
}
