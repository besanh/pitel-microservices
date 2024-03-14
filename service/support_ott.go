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
		authInfo.Source = agent.Source
		return authInfo, nil
	} else {
		filter := model.AgentAllocationFilter{
			ConversationId: newConversationId,
			MainAllocate:   "active",
		}
		total, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return authInfo, err
		}
		if total > 0 {
			authInfo.TenantId = (*agentAllocations)[0].TenantId
			authInfo.UserId = (*agentAllocations)[0].AgentId
			return authInfo, nil
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
						ConversationId: newConversationId,
						QueueId:        (*connectionQueues)[0].QueueId,
						MainAllocate:   "active",
					}
					totalAgentAllocation, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filterAgentAllocation, -1, 0)
					if err != nil {
						log.Error(err)
						return authInfo, err
					}
					if totalAgentAllocation > 0 {
						authInfo.TenantId = (*agentAllocations)[0].TenantId
						authInfo.UserId = (*agentAllocations)[0].AgentId
						authInfo.Source = (*agentAllocations)[0].Source
						for s := range WsSubscribers.Subscribers {
							if s.UserId == authInfo.UserId && (s.Level == "user" || s.Level == "agent") {
								agent = *s
							}
						}
						if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+newConversationId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
							log.Error(err)
							return authInfo, err
						}
						return authInfo, nil
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
									authInfo.Source = agent.Source
								}
							} else if strings.ToLower(chatRouting.RoutingAlias) == "round_robin_online" {
								agentTmp, err := RoundRobinAgentOnline(ctx, GenerateConversationId(message.AppId, message.ExternalUserId), queueAgents)
								if err != nil {
									log.Error(err)
									return authInfo, err
								}
								userLives = append(userLives, *agentTmp)
								agent = *agentTmp
								authInfo.TenantId = agent.TenantId
								authInfo.UserId = agent.UserId
								authInfo.Source = agent.Source
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
									MainAllocate:       "active",
									Source:             agent.Source,
								}
								log.Infof("conversation %s allocated to agent %s", newConversationId, agent.Username)
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
				log.Error("connect for conversation " + newConversationId + " not found")
				return authInfo, errors.New("connect for conversation " + newConversationId + " not found")
			}
		}
	}
}

func UpSertConversation(ctx context.Context, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
	newConversationId := GenerateConversationId(data.AppId, data.ExternalUserId)
	conversation = model.Conversation{
		TenantId:         data.TenantId,
		ConversationId:   newConversationId,
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
			ConversationId: []string{newConversationId},
			AppId:          []string{data.AppId},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, data.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
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
			return conversation, isNew, nil
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
	id = GenerateConversationId(conversation.AppId, conversation.ExternalUserId)
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
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, id)
	if err != nil {
		log.Error(err)
		return id, err
	} else if len(conversationExist.ExternalUserId) > 0 {
		log.Errorf("conversation %s not found", id)
		return id, errors.New("conversation " + id + " not found")
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, id, esDoc); err != nil {
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
