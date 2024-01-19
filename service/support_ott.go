package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func CheckChatQueueSetting(ctx context.Context, filter model.QueueFilter, externalUserId string) (string, error) {
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
		if len(WsSubscribers.Subscribers) > 0 {
			rand.NewSource(time.Now().UnixNano())
			randomIndex := rand.Intn(len(WsSubscribers.Subscribers))
			for s := range WsSubscribers.Subscribers {
				subscribers = append(subscribers, *s)
			}
			agent := subscribers[randomIndex]
			agentAllocationCache := cache.RCache.Get(AGENT_ALLOCATION + "_" + externalUserId)
			if agentAllocationCache != nil {
				agentTmp := Subscriber{}
				if err := json.Unmarshal([]byte(agentAllocationCache.(string)), &agentTmp); err != nil {
					log.Error(err)
					return agentId, err
				}
				agent = agentTmp
			} else {
				filter := model.AgentAllocationFilter{
					ConversationId: externalUserId,
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
				} else {
					agentId = agent.UserId
					agentAllocation := model.AgentAllocation{
						Base:               model.InitBase(),
						ConversationId:     externalUserId,
						AgentId:            agent.UserId,
						QueueId:            chatQueue.Id,
						AllocatedTimestamp: time.Now().Unix(),
					}
					if err := repository.AgentAllocationRepo.Insert(ctx, repository.DBConn, agentAllocation); err != nil {
						log.Error(err)
						return agentId, err
					}
				}

				if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+externalUserId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
					log.Error(err)
					return agentId, err
				}
			}

			if len(agent.UserId) > 0 {
				agentId = agent.UserId
				filter := model.AgentAllocationFilter{
					ConversationId: externalUserId,
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
		ConversationId:   data.ExternalUserId,
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		UserIdByApp:      data.UserIdByApp,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		ExternalUserId:   data.ExternalUserId,
		CreatedAt:        time.Now().Format(time.RFC3339),
	}

	isExisted := false
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.ExternalUserId)
	if conversationCache != nil {
		isExisted = true
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if err := UpdateESAndCache(ctx, data.AppId, conversation.ConversationId, conversation.ExternalUserId); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		return conversation, isNew, nil
	} else {
		filter := model.ConversationFilter{
			AppId:          []string{data.AppId},
			ConversationId: []string{data.ExternalUserId},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, data.AppId, ES_INDEX_CONVERSATION, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if total > 0 {
			conversation = (*conversations)[0]
			if err := cache.RCache.Set(CONVERSATION+"_"+data.ExternalUserId, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
			isExisted = true
			return conversation, isExisted, nil
		}
	}

	if !isExisted {
		id, err := InsertConversation(ctx, conversation)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		conversation.ConversationId = id
		if err := cache.RCache.Set(CONVERSATION+"_"+data.ExternalUserId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isNew = true
	}

	return conversation, isNew, nil
}

func InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error) {
	id = conversation.ExternalUserId
	// id = uuid.NewString()
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

func CheckConversationInAgent(userId string, allocationAgent []*model.AgentAllocation) bool {
	for _, item := range allocationAgent {
		if item.ConversationId == userId {
			return true
		}
	}
	return false
}

/**
* Update ES and Cache
* API get conversation can get from redis, and here can caching to descrese the number of api calls to ES
 */
func UpdateESAndCache(ctx context.Context, appId, conversationId, userId string) error {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, appId, ES_INDEX_CONVERSATION, userId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ExternalUserId) < 1 {
		log.Errorf("conversation %s not found", conversationId)
		return errors.New("conversation not found")
	}

	conversationExist.UpdatedAt = time.Now().Format(time.RFC3339)
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		log.Error(err)
		return err
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return err
	}
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, userId, esDoc); err != nil {
		log.Error(err)
		return err
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+userId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
