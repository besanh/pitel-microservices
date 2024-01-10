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
	queue := model.ChatQueue{}
	_, _, err := repository.ChatQueueRepo.GetQueue(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return agentId, err
	}
	// Get chat queue from cache or db
	ok, err := cache.RCache.IsHExisted(CHAT_QUEUE, filter.AppId)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else if !ok {
		total, queues, err := repository.ChatQueueRepo.GetQueue(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if total < 1 {
			log.Error("queue not found")
			return agentId, errors.New("queue not found")
		}
		queue = (*queues)[0]
		jsonByte, err := json.Marshal(&queue)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if err := cache.RCache.HSetRaw(ctx, CHAT_QUEUE, filter.AppId, string(jsonByte)); err != nil {
			log.Error(err)
			return agentId, err
		}
	}
	queuesCache, err := cache.RCache.HGet(CHAT_QUEUE, filter.AppId)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else {
		if err = json.Unmarshal([]byte(queuesCache), &queue); err != nil {
			log.Error(err)
			return agentId, err
		}
	}

	routing := model.ChatRouting{}
	// Get routing from cache or repository.DBConn
	ok, err = cache.RCache.IsHExisted(CHAT_ROUTING, queue.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else if !ok {
		routing, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, queue.ChatRoutingId)
		if err != nil {
			log.Error(err)
			return agentId, err
		} else if routing == nil {
			log.Error("routing not found")
			return agentId, errors.New("routing not found")
		}
		jsonByte, err := json.Marshal(&routing)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if err := cache.RCache.HSetRaw(ctx, CHAT_ROUTING, queue.ChatRoutingId, string(jsonByte)); err != nil {
			log.Error(err)
			return agentId, err
		}
	}
	routingCache, err := cache.RCache.HGet(CHAT_ROUTING, queue.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return agentId, err
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
	ok, err = cache.RCache.IsHExisted(CHAT_QUEUE_AGENT, queue.Id)
	if err != nil {
		log.Error(err)
		return agentId, err
	} else if !ok {
		total, agentDatas, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, repository.DBConn, filterChatQueueAgent, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		// if total < 1 {
		// 	log.Error("agent not found")
		// 	return agentId, errors.New("agent not found")
		// }
		agents = (*agentDatas)
		if total > 0 {
			for _, item := range *agentDatas {
				jsonByte, err := json.Marshal(&item)
				if err != nil {
					log.Error(err)
					return agentId, err
				}
				if err := cache.RCache.HSetRaw(ctx, CHAT_QUEUE_AGENT, item.AgentId, string(jsonByte)); err != nil {
					log.Error(err)
					return agentId, err
				}
			}
		}
	} else {
		chatQueueAgentCache, err := cache.RCache.HGetAll(CHAT_QUEUE_AGENT)
		if err != nil {
			log.Error(err)
			return agentId, err
		} else {
			if err := util.ParseAnyToAny(chatQueueAgentCache, &agents); err != nil {
				log.Error(err)
				return agentId, err
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

func GetConversationExist(ctx context.Context, data model.OttMessage) (conversation model.Conversation, isExisted bool, err error) {
	conversation = model.Conversation{
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		UserIdByApp:      data.UserIdByApp,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		Uid:              data.UserId,
	}

	ok, err := cache.RCache.IsHExisted(CONVERSATION, data.UserIdByApp)
	if err != nil {
		log.Error(err)
		return conversation, true, err
	} else if !ok {
		jsonByte, err := json.Marshal(&conversation)
		if err != nil {
			log.Error(err)
			return conversation, false, err
		}
		if err := cache.RCache.HSetRaw(ctx, CONVERSATION, data.UserIdByApp, string(jsonByte)); err != nil {
			log.Error(err)
			return conversation, false, err
		}
	}

	conversationCache, err := cache.RCache.HGet(CONVERSATION, data.UserIdByApp)
	if err != nil {
		log.Error(err)
		return conversation, false, err
	} else {
		err := json.Unmarshal([]byte(conversationCache), &conversation)
		if err != nil {
			log.Error(err)
			return conversation, false, err
		}
	}

	return conversation, true, nil
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
