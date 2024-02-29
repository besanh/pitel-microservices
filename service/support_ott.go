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
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

func CheckChatSetting(ctx context.Context, message model.Message) (model.AuthUser, error) {
	var authInfo model.AuthUser
	var userLives []Subscriber
	var agent Subscriber

	newConversationId := GenerateConversationId(message.AppId, message.ExternalUserId)
	agentAllocationCache := cache.RCache.Get(AGENT_ALLOCATION + "_" + newConversationId)
	if agentAllocationCache != nil {
		agentTmp := Subscriber{}
		if err := json.Unmarshal([]byte(agentAllocationCache.(string)), &agentTmp); err != nil {
			log.Error(err)
			return authInfo, err
		}
		agent = agentTmp
		authInfo.TenantId = agent.TenantId
		authInfo.UserId = agent.UserId
	} else {
		filter := model.AgentAllocationFilter{
			ConversationId: message.ExternalUserId,
		}
		total, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return authInfo, err
		}
		if total > 0 {
			authInfo.TenantId = (*agentAllocations)[0].TenantId
			authInfo.UserId = (*agentAllocations)[0].AgentId
		} else {
			// Get connection
			connectionType := ""
			if message.MessageType == "zalo" {
				connectionType = "zalo"
			} else if message.MessageType == "face" {
				connectionType = "facebook"
			}
			connectionFilter := model.ChatConnectionAppFilter{
				ConnectionType: connectionType,
				OaId:           message.OaId,
				AppId:          message.AppId,
				Status:         "active",
			}
			totalConnection, connectionApps, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, connectionFilter, 1, 0)
			if err != nil {
				log.Error(err)
				return authInfo, err
			}
			if totalConnection > 0 {
				filter := model.ConnectionQueueFilter{
					ConnectionId: (*connectionApps)[0].Id,
				}
				totalConnectionQueue, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, -1, 0)
				if err != nil {
					log.Error(err)
					return authInfo, err
				}
				if totalConnectionQueue > 0 {
					filterAgentAllocation := model.AgentAllocationFilter{
						ConversationId: message.ExternalUserId,
						QueueId:        (*connectionQueues)[0].QueueId,
					}
					totalAgentAllocation, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filterAgentAllocation, -1, 0)
					if err != nil {
						log.Error(err)
						return authInfo, err
					}
					if totalAgentAllocation > 0 {
						authInfo.TenantId = (*agentAllocations)[0].TenantId
						authInfo.UserId = (*agentAllocations)[0].AgentId
						for s := range WsSubscribers.Subscribers {
							if s.UserId == authInfo.UserId && (s.Level == "user" || s.Level == "agent") {
								agent = *s
							}
						}
						if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+message.AppId+"_"+message.ExternalUserId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
							log.Error(err)
							return authInfo, err
						}
					} else {
						// Connection prevent duplicate
						// Meaning: 1 connection with page A in 1 app => only recieve one queue
						queue, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, (*connectionQueues)[0].QueueId)
						if err != nil {
							log.Error(err)
							return authInfo, err
						} else if queue == nil {
							log.Error("queue " + (*connectionQueues)[0].QueueId + " not found")
							return authInfo, errors.New("queue " + (*connectionQueues)[0].QueueId + " not found")
						}

						chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, queue.ChatRoutingId)
						if err != nil {
							log.Error(err)
							return authInfo, err
						} else if chatRouting == nil {
							log.Error("chat routing " + queue.ChatRoutingId + " not found")
							return authInfo, errors.New("chat routing " + queue.ChatRoutingId + " not found")
						}

						filterQueueAgent := model.ChatQueueAgentFilter{
							QueueId: []string{queue.Id},
						}
						totalQueueAgents, queueAgents, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, repository.DBConn, filterQueueAgent, -1, 0)
						if err != nil {
							log.Error(err)
							return authInfo, err
						}
						if totalQueueAgents > 0 {
							if strings.ToLower(chatRouting.RoutingAlias) == "random" {
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
								if len(userLives) > 0 {
									rand.NewSource(time.Now().UnixNano())
									randomIndex := rand.Intn(len(userLives))
									agent = userLives[randomIndex]
									authInfo.TenantId = agent.TenantId
									authInfo.UserId = agent.UserId
								}
							} else if strings.ToLower(chatRouting.RoutingAlias) == "round_robin_online" {
								agentTmp, err := RoundRobinAgentOnline(ctx, GenerateConversationId(message.AppId, message.ExternalUserId))
								if err != nil {
									log.Error(err)
									return authInfo, err
								}
								userLives = append(userLives, *agentTmp)
								agent = *agentTmp
								authInfo.TenantId = agent.TenantId
								authInfo.UserId = agent.UserId
							}

							if len(userLives) > 0 {
								newConversationId = GenerateConversationId(message.AppId, message.ExternalUserId)
								agentAllocation := model.AgentAllocation{
									Base:               model.InitBase(),
									TenantId:           (*connectionApps)[0].TenantId,
									ConversationId:     newConversationId,
									AppId:              message.AppId,
									AgentId:            agent.UserId,
									QueueId:            queue.Id,
									AllocatedTimestamp: time.Now().Unix(),
								}
								log.Info(agent.Username)
								if err := repository.AgentAllocationRepo.Insert(ctx, repository.DBConn, agentAllocation); err != nil {
									log.Error(err)
									return authInfo, err
								}

								if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+newConversationId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
									log.Error(err)
									return authInfo, err
								}

								return authInfo, nil
							} else {
								log.Error("agent not available")
								return authInfo, errors.New("agent not available")
							}
						} else {
							log.Error("queue agent not found")
							return authInfo, errors.New("queue agent not found")
						}
					}
				} else {
					log.Error("queue not found")
					return authInfo, errors.New("queue not found")
				}
			} else {
				log.Error("connection not found")
				return authInfo, errors.New("connection " + newConversationId + " not found")
			}
		}
	}
	return authInfo, nil
}

func UpSertConversation(ctx context.Context, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
	conversation = model.Conversation{
		TenantId:         data.TenantId,
		ConversationId:   data.ExternalUserId,
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		ShareInfo:        data.ShareInfo,
		ExternalUserId:   data.ExternalUserId,
		CreatedAt:        time.Now().Format(time.RFC3339),
	}

	isExisted := false
	newConversationId := GenerateConversationId(data.AppId, data.ExternalUserId)
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + newConversationId)
	if conversationCache != nil {
		isExisted = true
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if err := UpdateESAndCache(ctx, data.TenantId, data.AppId, data.ExternalUserId, *conversation.ShareInfo); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		return conversation, isNew, nil
	} else {
		filter := model.ConversationFilter{
			ConversationId: []string{data.ExternalUserId},
			AppId:          []string{data.AppId},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if total > 0 {
			if err := util.ParseAnyToAny((*conversations)[0], &conversation); err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
			conversation.ShareInfo = data.ShareInfo

			tmpBytes, err := json.Marshal(conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
			esDoc := map[string]any{}
			if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
			newConversationId := GenerateConversationId(conversation.AppId, conversation.ExternalUserId)
			if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversation.AppId, newConversationId, esDoc); err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
			if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
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
		if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isNew = true
	}

	return conversation, isNew, nil
}

func InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error) {
	id = conversation.ExternalUserId
	newConversationId := GenerateConversationId(conversation.AppId, id)
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
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
		log.Error(err)
		return id, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
			log.Error(err)
			return id, err
		}
	}
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, newConversationId)
	if err != nil {
		log.Error(err)
		return id, err
	} else if len(conversationExist.ExternalUserId) > 0 {
		log.Errorf("conversation %s not found", id)
		return id, errors.New("conversation " + id + " not found")
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, newConversationId, esDoc); err != nil {
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
func UpdateESAndCache(ctx context.Context, tenantId, appId, conversationId string, shareInfo model.ShareInfo) error {
	newConversationId := GenerateConversationId(appId, conversationId)
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ExternalUserId) < 1 {
		log.Errorf("conversation %s not found", newConversationId)
		return errors.New("conversation " + newConversationId + " not found")
	}

	conversationExist.ShareInfo = &shareInfo
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
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
		log.Error(err)
		return err
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func GenerateConversationId(appId, conversationId string) (newConversationId string) {
	newConversationId = appId + "_" + conversationId
	return
}

func RoundRobinAgentOnline(ctx context.Context, conversationId string) (*Subscriber, error) {
	userLive := Subscriber{}
	userLives := []Subscriber{}
	for s := range WsSubscribers.Subscribers {
		if s.Level == "user" || s.Level == "agent" {
			userLives = append(userLives, *s)
		}
	}
	if len(userLives) > 0 {
		isOk, index, userAllocatePrevious := GetAgentIsRoundRobin(userLives, conversationId)
		if isOk {
			if (index+1)%len(userLives) <= len(userLives) {
				userLive = userLives[(index+1)%len(userLives)]
			} else {
				userLive = userLives[0]
			}
		} else {
			userLive = *userAllocatePrevious
		}
		if err := cache.RCache.Set(AGENT_ROUND_ROBIN_ONLINE+"_"+conversationId+"_"+userLive.Id, userLive, AGENT_ROUND_ROBIN_ONLINE_EXPIRE); err != nil {
			log.Error(err)
			return &userLive, err
		}
	} else {
		return &userLive, errors.New("no user online")
	}
	return &userLive, nil
}

func GetAgentIsRoundRobin(userLives []Subscriber, conversationId string) (bool, int, *Subscriber) {
	isOk := false
	index := 0
	userLive := Subscriber{}
	for i, item := range userLives {
		id := conversationId + "_" + item.Id
		itemCache := cache.RCache.Get(AGENT_ROUND_ROBIN_ONLINE + "_" + id)
		if itemCache != nil {
			return isOk, index, &userLives[0]
		} else if itemCache == nil {
			return isOk, index, &userLives[0]
		}
		userLiveTmp := Subscriber{}
		if err := util.ParseAnyToAny(itemCache, &userLiveTmp); err != nil {
			log.Error(err)
			return isOk, index, &userLives[0]
		}
		isOk = true
		index = i
		userLive = userLiveTmp
		break
	}
	if isOk {
		return isOk, index, &userLive
	}
	userLive = userLives[0]
	return isOk, index, &userLive
}
