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
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func CheckChatQueueSetting(ctx context.Context, filter model.QueueFilter) (string, error) {
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
		routing, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, chatQueue.ChatRoutingId)
		if err != nil {
			log.Error(err)
			return agentId, err
		} else if routing == nil {
			log.Error("routing not found")
			return agentId, errors.New("routing not found")
		}
		if err := cache.RCache.Set(CHAT_ROUTING+"_"+chatQueue.ChatRoutingId, routing, CHAT_ROUTING_EXPIRE); err != nil {
			log.Error(err)
			return agentId, err
		}
	}

	agents := []model.ChatQueueAgent{}
	// Chat queue agent
	filterChatQueueAgent := model.ChatQueueAgentFilter{
		QueueId: chatQueue.Id,
	}
	chatQueueAgentCache := cache.RCache.Get(CHAT_QUEUE_AGENT + "_" + chatQueue.Id)
	if chatQueueAgentCache != nil {
		if err := json.Unmarshal([]byte(chatQueueAgentCache.(string)), &agents); err != nil {
			log.Error(err)
			return agentId, err
		}
	} else {
		total, agentDatas, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, repository.DBConn, filterChatQueueAgent, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		agents = (*agentDatas)
		if total > 0 {
			for _, item := range *agentDatas {
				if err := cache.RCache.Set(CHAT_QUEUE_AGENT+"_"+item.AgentId, item, CHAT_QUEUE_AGENT_EXPIRE); err != nil {
					log.Error(err)
					return agentId, err
				}
			}
		}
	}

	if routing.RoutingName == "random" {
		if len(agents) > 0 {
			rand.NewSource(time.Now().UnixNano())
			randomIndex := rand.Intn(len(agents))
			agent := agents[randomIndex]
			agentId = agent.AgentId
		}
	} else if routing.RoutingName == "min_conversation" {
	}

	return agentId, nil
}

func GetConversationExist(ctx context.Context, data model.OttMessage) (conversation model.Conversation, err error) {
	conversation = model.Conversation{
		ConversationId:   uuid.NewString(),
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		UserIdByApp:      data.UserIdByApp,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		Uid:              data.UserId,
	}

	isExisted := false
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.UserIdByApp)
	if conversationCache != nil {
		isExisted = true
		if err := util.ParseAnyToAny(conversationCache, &conversation); err != nil {
			log.Error(err)
			return conversation, err
		}
		return conversation, nil
	} else {
		filter := model.ConversationFilter{
			AppId:       []string{data.AppId},
			UserIdByApp: []string{data.UserIdByApp},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, data.AppId, ES_INDEX_CONVERSATION, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return conversation, err
		}
		if total > 0 {
			isExisted = true
			conversation = (*conversations)[0]
			if err := cache.RCache.Set(CONVERSATION+"_"+data.UserIdByApp, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return conversation, err
			}
		}
	}

	if !isExisted {
		id, err := InsertConversation(ctx, conversation)
		if err != nil {
			log.Error(err)
			return conversation, err
		}
		conversation.ConversationId = id
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
