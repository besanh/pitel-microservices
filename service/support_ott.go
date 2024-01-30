package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

func CheckChatSetting(ctx context.Context, message model.Message) (string, error) {
	var agentId string
	var userLives []Subscriber
	var agent Subscriber

	agentAllocationCache := cache.RCache.Get(AGENT_ALLOCATION + "_" + message.ExternalUserId)
	if agentAllocationCache != nil {
		agentTmp := Subscriber{}
		if err := json.Unmarshal([]byte(agentAllocationCache.(string)), &agentTmp); err != nil {
			log.Error(err)
			return agentId, err
		}
		agent = agentTmp
		agentId = agent.UserId
	} else {
		filter := model.AgentAllocationFilter{
			ConversationId: message.ExternalUserId,
		}
		total, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return agentId, err
		}
		if total > 0 {
			agentId = (*agentAllocations)[0].AgentId
		} else {
			// Get connection
			connectionFilter := model.ChatConnectionAppFilter{
				OaId:   message.OaId,
				AppId:  message.AppId,
				Status: "active",
			}
			totalConnection, connectionApps, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, connectionFilter, 1, 0)
			if err != nil {
				log.Error(err)
				return agentId, err
			}
			if totalConnection > 0 {
				filter := model.ConnectionQueueFilter{
					ConnectionId: (*connectionApps)[0].Id,
				}
				totalConnectionQueue, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, -1, 0)
				if err != nil {
					log.Error(err)
					return agentId, err
				}
				if totalConnectionQueue > 0 {
					filterAgentAllocation := model.AgentAllocationFilter{
						ConversationId: message.ExternalUserId,
						QueueId:        (*connectionQueues)[0].QueueId,
					}
					totalAgentAllocation, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filterAgentAllocation, -1, 0)
					if err != nil {
						log.Error(err)
						return agentId, err
					}
					if totalAgentAllocation > 0 {
						agentId = (*agentAllocations)[0].AgentId
						for s := range WsSubscribers.Subscribers {
							if s.UserId == agentId && (s.Level == "user" || s.Level == "agent") {
								agent = *s
							}
						}
						if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+message.ExternalUserId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
							log.Error(err)
							return agentId, err
						}
					} else {
						// Connection prevent duplicate
						// Meaning: 1 connection with page A in 1 app => only recieve one queue
						queue, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, (*connectionQueues)[0].QueueId)
						if err != nil {
							log.Error(err)
							return agentId, err
						} else if queue == nil {
							log.Error("queue " + (*connectionQueues)[0].QueueId + " not found")
							return agentId, errors.New("queue " + (*connectionQueues)[0].QueueId + " not found")
						}

						chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, queue.ChatRoutingId)
						if err != nil {
							log.Error(err)
							return agentId, err
						} else if chatRouting == nil {
							log.Error("chat routing " + queue.ChatRoutingId + " not found")
							return agentId, errors.New("chat routing " + queue.ChatRoutingId + " not found")
						}

						filterQueueAgent := model.ChatQueueAgentFilter{
							QueueId: []string{queue.Id},
						}
						totalQueueAgents, queueAgents, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, repository.DBConn, filterQueueAgent, -1, 0)
						if err != nil {
							log.Error(err)
							return agentId, err
						}
						if totalQueueAgents > 0 {
							if strings.ToLower(chatRouting.RoutingName) == "random" {
								agentUuids := []string{}
								for _, item := range *queueAgents {
									agentUuids = append(agentUuids, item.AgentId)
								}
								for s := range WsSubscribers.Subscribers {
									if slices.Contains[[]string](agentUuids, s.UserId) && s.Level == "user" || s.Level == "agent" {
										userLives = append(userLives, *s)
									}
								}

								// Pick random
								rand.NewSource(time.Now().UnixNano())
								randomIndex := rand.Intn(len(userLives))
								agent = userLives[randomIndex]
								agentId = agent.UserId
							} else {
							}

							if len(userLives) > 0 {
								agentAllocation := model.AgentAllocation{
									Base:               model.InitBase(),
									ConversationId:     message.ExternalUserId,
									AgentId:            agent.UserId,
									QueueId:            queue.Id,
									AllocatedTimestamp: time.Now().Unix(),
								}
								if err := repository.AgentAllocationRepo.Insert(ctx, repository.DBConn, agentAllocation); err != nil {
									log.Error(err)
									return agentId, err
								}

								if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+message.ExternalUserId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
									log.Error(err)
									return agentId, err
								}

								return agentId, nil
							} else {
								log.Error("agent not available")
								return agentId, errors.New("agent not available")
							}
						} else {
							log.Error("queue agent not found")
							return agentId, errors.New("queue agent not found")
						}
					}
				} else {
					log.Error("queue not found")
					return agentId, errors.New("queue not found")
				}
			} else {
				log.Error("connection not found")
				return agentId, errors.New("connection " + message.OaId + " not found")
			}
		}
	}
	return agentId, nil
}

func UpSertConversation(ctx context.Context, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
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
		if err := UpdateESAndCache(ctx, data.AppId, data.ExternalUserId); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		return conversation, isNew, nil
	} else {
		filter := model.ConversationFilter{
			ConversationId: []string{data.ExternalUserId},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, 1, 0)
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
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, conversation.AppId, ES_INDEX_CONVERSATION, id)
	if err != nil {
		log.Error(err)
		return id, err
	} else if len(conversationExist.ExternalUserId) > 0 {
		log.Errorf("conversation %s not found", id)
		return id, errors.New("conversation not found")
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
func UpdateESAndCache(ctx context.Context, appId, conversationId string) error {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, appId, ES_INDEX_CONVERSATION, conversationId)
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
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversationId, esDoc); err != nil {
		log.Error(err)
		return err
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+conversationId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
