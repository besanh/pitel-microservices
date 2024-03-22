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

func CheckChatSetting(ctx context.Context, message model.Message) (model.User, error) {
	var user model.User
	var authInfo model.AuthUser
	var userLives []Subscriber
	var agent model.AgentAllocate
	var isOk bool

	newConversationId := GenerateConversationId(message.AppId, message.ExternalUserId)
	agentAllocationCache := cache.RCache.Get(AGENT_ALLOCATION + "_" + newConversationId)
	if agentAllocationCache != nil {
		agentTmp := model.AgentAllocate{}
		if err := json.Unmarshal([]byte(agentAllocationCache.(string)), &agentTmp); err != nil {
			log.Error(err)
			user.AuthUser = &authInfo
			user.IsOk = isOk
			return user, err
		}
		agent = agentTmp
		authInfo.TenantId = agent.TenantId
		authInfo.UserId = agent.AgentId
		authInfo.Username = agent.Username
		user.AuthUser = &authInfo
		user.IsOk = true
		user.ConnectionId = agent.ConnectionId
		user.QueueId = agent.QueueId

		return user, nil
	} else {
		filter := model.AgentAllocateFilter{
			ConversationId: newConversationId,
			MainAllocate:   "active",
		}
		total, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return user, err
		}
		if total > 0 {
			authInfo.TenantId = (*agentAllocations)[0].TenantId
			authInfo.UserId = (*agentAllocations)[0].AgentId
			user.AuthUser = &authInfo
			user.IsOk = true
			user.ConnectionId = (*agentAllocations)[0].ConnectionId
			user.QueueId = (*agentAllocations)[0].QueueId

			if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+newConversationId, (*agentAllocations)[0], AGENT_ALLOCATION_EXPIRE); err != nil {
				log.Error(err)
				return user, err
			}

			return user, nil
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
				return user, err
			}
			if totalConnection > 0 {
				filter := model.ConnectionQueueFilter{
					ConnectionId: (*connectionApps)[0].Id,
				}
				totalConnectionQueue, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, 1, 0)
				if err != nil {
					log.Error(err)
					return user, err
				}
				if totalConnectionQueue > 0 {
					filterAgentAllocation := model.AgentAllocateFilter{
						ConversationId: newConversationId,
						QueueId:        (*connectionQueues)[0].QueueId,
						MainAllocate:   "active",
					}
					totalAgentAllocation, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filterAgentAllocation, -1, 0)
					if err != nil {
						log.Error(err)
						return user, err
					}
					if totalAgentAllocation > 0 {
						authInfo.TenantId = (*agentAllocations)[0].TenantId
						authInfo.UserId = (*agentAllocations)[0].AgentId

						for s := range WsSubscribers.Subscribers {
							if s.UserId == authInfo.UserId && (s.Level == "user" || s.Level == "agent") {
								agent.AgentId = s.UserId
							}
						}
						if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+newConversationId, agent, AGENT_ALLOCATION_EXPIRE); err != nil {
							log.Error(err)
							return user, err
						}
						user.IsOk = true
						user.AuthUser = &authInfo
						user.ConnectionId = (*agentAllocations)[0].ConnectionId
						user.QueueId = (*agentAllocations)[0].QueueId

						return user, nil
					} else {
						// Connection prevent duplicate
						// Meaning: 1 connection with page A in 1 app => only recieve one queue
						queue, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, (*connectionQueues)[0].QueueId)
						if err != nil {
							log.Error(err)
							return user, err
						} else if queue == nil {
							log.Error("queue " + (*connectionQueues)[0].QueueId + " not found")
							return user, errors.New("queue " + (*connectionQueues)[0].QueueId + " not found")
						}

						chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, queue.ChatRoutingId)
						if err != nil {
							log.Error(err)
							return user, err
						} else if chatRouting == nil {
							log.Error("chat routing " + queue.ChatRoutingId + " not found")
							return user, errors.New("chat routing " + queue.ChatRoutingId + " not found")
						}

						filterQueueAgent := model.ChatQueueAgentFilter{
							QueueId: []string{queue.Id},
						}
						totalQueueAgents, queueAgents, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, repository.DBConn, filterQueueAgent, -1, 0)
						if err != nil {
							log.Error(err)
							return user, err
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
									tmp := userLives[randomIndex]
									agent.TenantId = tmp.TenantId
									agent.AgentId = tmp.UserId
									agent.Username = tmp.Username

									authInfo.TenantId = agent.TenantId
									authInfo.UserId = agent.AgentId
								}
							} else if strings.ToLower(chatRouting.RoutingAlias) == "round_robin_online" {
								agentTmp, err := RoundRobinAgentOnline(ctx, GenerateConversationId(message.AppId, message.ExternalUserId), queueAgents)
								if err != nil {
									log.Error(err)
									return user, err
								}
								userLives = append(userLives, *agentTmp)
								agent.TenantId = agentTmp.TenantId
								agent.AgentId = agentTmp.UserId
								agent.Username = agentTmp.Username

								authInfo.TenantId = agent.TenantId
							}

							if len(userLives) > 0 {
								newConversationId = GenerateConversationId(message.AppId, message.ExternalUserId)
								agentAllocation := model.AgentAllocate{
									Base:               model.InitBase(),
									TenantId:           (*connectionApps)[0].TenantId,
									ConversationId:     newConversationId,
									AppId:              message.AppId,
									AgentId:            agent.AgentId,
									QueueId:            queue.Id,
									AllocatedTimestamp: time.Now().Unix(),
									MainAllocate:       "active",
									ConnectionId:       (*connectionQueues)[0].ConnectionId,
								}
								log.Infof("conversation %s allocated to agent %s", newConversationId, agent.Username)
								if err := repository.AgentAllocationRepo.Insert(ctx, repository.DBConn, agentAllocation); err != nil {
									log.Error(err)
									return user, err
								}

								if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+newConversationId, agentAllocation, AGENT_ALLOCATION_EXPIRE); err != nil {
									log.Error(err)
									return user, err
								}

								user.IsOk = true
								user.AuthUser = &authInfo
								user.ConnectionId = (*connectionQueues)[0].ConnectionId
								user.QueueId = (*connectionQueues)[0].QueueId

								return user, nil
							} else {
								log.Error("agent not available")
								user.IsOk = true
								return user, errors.New("agent not available")
							}
						} else {
							log.Error("queue agent not found")
							return user, errors.New("queue agent not found")
						}
					}
				} else {
					log.Error("queue not found")
					return user, errors.New("queue not found")
				}
			} else {
				log.Error("connect for conversation " + newConversationId + " not found")
				return user, errors.New("connect for conversation " + newConversationId + " not found")
			}
		}
	}
}

func UpSertConversation(ctx context.Context, connectionId string, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
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
	shareInfo := data.ShareInfo

	isExisted := false
	// conversationCache := cache.RCache.Get(CONVERSATION + "_" + newConversationId)
	// if conversationCache != nil {
	// 	isExisted = true
	// 	if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
	// 		log.Error(err)
	// 		return conversation, isNew, err
	// 	}
	// 	if err := UpdateESAndCache(ctx, data.TenantId, data.AppId, data.ExternalUserId, connectionId, *conversation.ShareInfo); err != nil {
	// 		log.Error(err)
	// 		return conversation, isNew, err
	// 	}
	// 	return conversation, isNew, nil
	// } else {
	log.Info("zz", data.ShareInfo)
	filter := model.ConversationFilter{
		ConversationId: []string{newConversationId},
		AppId:          []string{data.AppId},
	}
	total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return conversation, isNew, err
	}
	if total > 0 {
		conversation.TenantId = (*conversations)[0].TenantId
		conversation.ConversationId = (*conversations)[0].ConversationId
		conversation.ConversationType = (*conversations)[0].ConversationType
		conversation.AppId = (*conversations)[0].AppId
		conversation.OaId = (*conversations)[0].OaId
		conversation.OaName = (*conversations)[0].OaName
		conversation.OaAvatar = (*conversations)[0].OaAvatar
		conversation.ExternalUserId = (*conversations)[0].ExternalUserId
		conversation.Username = (*conversations)[0].Username
		conversation.Username = (*conversations)[0].Username
		conversation.Avatar = (*conversations)[0].Avatar
		conversation.IsDone = (*conversations)[0].IsDone
		conversation.IsDoneBy = (*conversations)[0].IsDoneBy

		conversation.ShareInfo = shareInfo
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}

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
	// }

	if !isExisted {
		id, err := InsertConversation(ctx, conversation, connectionId)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		conversation.ConversationId = id
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}
		if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isNew = true
	}

	return conversation, isNew, nil
}

func InsertConversation(ctx context.Context, conversation model.Conversation, connectionId string) (id string, err error) {
	id = GenerateConversationId(conversation.AppId, conversation.ExternalUserId)
	if len(connectionId) > 0 {
		conversation, err = CacheConnection(ctx, connectionId, conversation)
		if err != nil {
			log.Error(err)
			return id, err
		}
	}
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
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, conversation.AppId, id)
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

func CheckConversationInAgent(userId string, allocationAgent []*model.AgentAllocate) bool {
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
func UpdateESAndCache(ctx context.Context, tenantId, appId, conversationId, connectionId string, shareInfo model.ShareInfo) error {
	var isUpdate bool
	newConversationId := GenerateConversationId(appId, conversationId)
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ExternalUserId) < 1 {
		isUpdate = true
		// Use when routing is pitel_bss_conversation_
		conversationExistSecond, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, appId, newConversationId)
		if err != nil {
			log.Error(err)
			return err
		} else if len(conversationExistSecond.ExternalUserId) < 1 {
			log.Errorf("conversation %s not found", newConversationId)
			return errors.New("conversation " + newConversationId + " not found")
		}
		if err := util.ParseAnyToAny(conversationExistSecond, &conversationExist); err != nil {
			log.Error(err)
			return err
		}
		conversationExist = conversationExistSecond
	}

	conversationExist.ShareInfo = &shareInfo
	conversationExist.UpdatedAt = time.Now().Format(time.RFC3339)
	conversationExist.TenantId = tenantId
	if len(connectionId) > 0 {
		*conversationExist, err = CacheConnection(ctx, connectionId, *conversationExist)
		if err != nil {
			log.Error(err)
			return err
		}
	}
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

	// PROBLEM: Use for conv have not authUser and then having authUser
	// TODO: insert new conv
	if isUpdate {
		if err := repository.ESRepo.DeleteById(ctx, ES_INDEX_CONVERSATION, newConversationId); err != nil {
			log.Error(err)
			return err
		}
		if err := repository.ESRepo.InsertLog(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return err
		}
	} else {
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return err
		}
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
